package main

import (
        "time"
	etcd "github.com/coreos/etcd/clientv3"
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoints": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"etcd_key": KeyResource(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	var endpoints []string
	values := d.Get("endpoints").([]interface{})
	for _, value := range values {
		endpoints = append(endpoints, value.(string))
	}
	config := etcd.Config{
		Endpoints: endpoints,
                DialTimeout: 5 * time.Second,
	}
	client, err := etcd.New(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}
