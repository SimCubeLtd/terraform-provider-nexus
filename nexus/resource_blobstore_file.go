/*
Use this resource to create a Nexus file blobstore.

Example Usage

```hcl
resource "nexus_blobstore_file" "default" {
  name = "blobstore-file"
  type = "File"
  path = "/nexus-data/blobstore-file"

  soft_quota {
    limit = 1024000000
    type  = "spaceRemainingQuota"
  }
}
```

*/
package nexus

import (
	"fmt"
	"log"

	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/blobstore"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceBlobstoreFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceBlobstoreFileCreate,
		Read:   resourceBlobstoreFileRead,
		Update: resourceBlobstoreFileUpdate,
		Delete: resourceBlobstoreFileDelete,
		Exists: resourceBlobstoreFileExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Blobstore name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"path": {
				Description: "The path to the blobstore contents. This can be an absolute path to anywhere on the system nxrm has access to or it can be a path relative to the sonatype-work directory",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"available_space_in_bytes": {
				Description: "Available space in Bytes",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"blob_count": {
				Description: "Count of blobs",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"soft_quota": {
				Description: "Soft quota of the blobstore",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limit": {
							Description:  "The limit in Bytes. Minimum value is 1000000",
							Required:     true,
							Type:         schema.TypeInt,
							ValidateFunc: validation.IntAtLeast(100000),
						},
						"type": {
							Description:  "The type to use such as spaceRemainingQuota, or spaceUsedQuota",
							Required:     true,
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"spaceRemainingQuota", "spaceUsedQuota"}, false),
						},
					},
				},
				MaxItems: 1,
				Optional: true,
				Type:     schema.TypeList,
			},
			"total_size_in_bytes": {
				Description: "The total size of the blobstore in Bytes",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func getBlobstoreFileFromResourceData(resourceData *schema.ResourceData) blobstore.File {
	bs := blobstore.File{
		Name: resourceData.Get("name").(string),
	}

	if _, ok := resourceData.GetOk("path"); ok {
		bs.Path = resourceData.Get("path").(string)
	}

	if _, ok := resourceData.GetOk("soft_quota"); ok {
		softQuotaConfig := resourceData.Get("soft_quota").([]interface{})[0].(map[string]interface{})

		bs.SoftQuota = &blobstore.SoftQuota{
			Limit: int64(softQuotaConfig["limit"].(int)),
			Type:  softQuotaConfig["type"].(string),
		}
	}

	return bs
}

func resourceBlobstoreFileCreate(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(*nexus.NexusClient)

	bs := getBlobstoreFileFromResourceData(resourceData)

	if err := nexusClient.BlobStore.File.Create(&bs); err != nil {
		return err
	}

	resourceData.SetId(bs.Name)
	err := resourceData.Set("name", bs.Name)
	if err != nil {
		return err
	}

	return resourceBlobstoreRead(resourceData, m)
}

func resourceBlobstoreFileRead(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(*nexus.NexusClient)

	bs, err := nexusClient.BlobStore.File.Get(resourceData.Id())
	log.Print(bs)
	if err != nil {
		return err
	}

	var genericBlobstoreInformation blobstore.Generic
	genericBlobstores, err := nexusClient.BlobStore.List()
	if err != nil {
		return err
	}
	for _, generic := range genericBlobstores {
		if generic.Name == bs.Name {
			genericBlobstoreInformation = generic
		}
	}

	if bs == nil {
		resourceData.SetId("")
		return nil
	}

	if err := resourceData.Set("available_space_in_bytes", genericBlobstoreInformation.AvailableSpaceInBytes); err != nil {
		return err
	}
	if err := resourceData.Set("blob_count", genericBlobstoreInformation.BlobCount); err != nil {
		return err
	}
	if err := resourceData.Set("name", bs.Name); err != nil {
		return err
	}
	if err := resourceData.Set("path", bs.Path); err != nil {
		return err
	}
	if err := resourceData.Set("total_size_in_bytes", genericBlobstoreInformation.TotalSizeInBytes); err != nil {
		return err
	}

	if bs.SoftQuota != nil {
		if err := resourceData.Set("soft_quota", flattenBlobstoreSoftQuota(bs.SoftQuota)); err != nil {
			return fmt.Errorf("error reading soft quota: %s", err)
		}
	}

	return nil
}

func resourceBlobstoreFileUpdate(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(*nexus.NexusClient)

	bs := getBlobstoreFileFromResourceData(resourceData)
	if err := nexusClient.BlobStore.File.Update(resourceData.Id(), &bs); err != nil {
		return err
	}

	return nil
}

func resourceBlobstoreFileDelete(resourceData *schema.ResourceData, m interface{}) error {
	nexusClient := m.(*nexus.NexusClient)

	if err := nexusClient.BlobStore.File.Delete(resourceData.Id()); err != nil {
		return err
	}

	resourceData.SetId("")

	return nil
}

func resourceBlobstoreFileExists(resourceData *schema.ResourceData, m interface{}) (bool, error) {
	nexusClient := m.(*nexus.NexusClient)

	bs, err := nexusClient.BlobStore.File.Get(resourceData.Id())
	return bs != nil, err
}