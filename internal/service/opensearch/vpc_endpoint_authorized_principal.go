// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0


package opensearch

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/opensearchservice"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// @SDKResource("aws_opensearch_vpc_endpoint_authorized_principal")
func ResourceVpcEndpointAuthorizedPrincipal() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceVpcEndpointAuthorizedPrincipalCreate,
		ReadWithoutTimeout:  resourceAwsOpenSearchDomainVpcEndpointAuthorizedPrincipalRead,
		UpdateWithoutTimeout: resourceAwsOpenSearchDomainVpcEndpointAuthorizedPrincipalUpdate,
		DeleteWithoutTimeout: resourceAwsOpenSearchDomainVpcEndpointAuthorizedPrincipalDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account": {
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceVpcEndpointAuthorizedPrincipalCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).OpenSearchConn(ctx)

	input := &opensearchservice.AuthorizeVpcEndpointAccessInput{
		DomainName: aws.String(d.Get("domain_name").(string)),
		Account: aws.String(d.Get("account").(string))
	}

	output, err := conn.AuthorizeVpcEndpointAccess(input)

	if err != nil{
		return sdkdiag.AppendErrorf(diags, "Error authorizing VPC endpoint access: %s", err)
	}

	d.SetId(aws.StringValue(output.Acount))

	if err := waitForVpcEndpointAuthorizedPrincipalCreated(conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	return append(diags)
}
