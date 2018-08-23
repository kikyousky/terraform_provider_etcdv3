package main

import (
	etcd "github.com/coreos/etcd/clientv3"
        "github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/hashicorp/terraform/helper/schema"
	"context"
        "time"
)

func KeyResource() *schema.Resource {
	return &schema.Resource{
		Create: createKey,
		Read:   readKey,
		Update: createKey,
		Delete: deleteKey,
		Schema: map[string]*schema.Schema{
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func createKey(d *schema.ResourceData, m interface{}) error {
	key := d.Get("key").(string)
	value := d.Get("value").(string)

        kv := m.(*etcd.Client)
        ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
	_, err := kv.Put(ctx, key, value)
        cancel()
	if err != nil {
		return err
	}

	d.SetId(key)
	d.Set("key", key)
	d.Set("value", value)

	return nil
}

func readKey(d *schema.ResourceData, m interface{}) error {
        kv := m.(*etcd.Client)

        ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
	_, err := kv.Get(ctx, d.Id())
        cancel()
	if err != nil {
                if err == rpctypes.ErrGRPCKeyNotFound {
			d.SetId("")
			return nil
                }
		return err
	}

	return nil
}

func deleteKey(d *schema.ResourceData, m interface{}) error {
        kv := m.(*etcd.Client)
        ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
	_, err := kv.Delete(ctx, d.Id())
        cancel()
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
