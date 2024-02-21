// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0


package opensearch

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/opensearchservice"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
)

// @SDKResource("aws_opensearch_authorized_principal")
func ResourceAuthorizedPrincipal() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAuthorizedPrincipalUpsert,
		ReadWithoutTimeout:  resourceAuthorizedPrincipalRead,
		UpdateWithoutTimeout: resourceAuthorizedPrincipalUpsert,
		DeleteWithoutTimeout: resourceAuthorizedPrincipalDelete,

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

func resourceAuthorizedPrincipalUpsert(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	
	domain_name := d.Get("domain_name").(string)

	conn := meta.(*conns.AWSClient).OpenSearchConn(ctx)

	input := &opensearchservice.AuthorizeVpcEndpointAccessInput{
		DomainName: aws.String(d.Get("domain_name").(string)),
		Account: aws.String(d.Get("account").(string)),
	}

	output, err := conn.AuthorizeVpcEndpointAccess(input)

	if err != nil{
		return sdkdiag.AppendErrorf(diags, "Error authorizing Principal %s", err)
	}

	d.SetId("authorized-principal-" + *output.AuthorizedPrincipal.Principal + "-" + *output.AuthorizedPrincipal.PrincipalType + "-" + domain_name)

	if err := waitForDomainUpdate(ctx, conn, domain_name, d.Timeout(schema.TimeoutCreate)); err != nil {
		return sdkdiag.AppendErrorf(diags, "Error authorizing principal %s: %s", d.Id(), err)
	}

	return append(diags, resourceAuthorizedPrincipalRead(ctx, d, meta)...)
}


func resourceAuthorizedPrincipalRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).OpenSearchConn(ctx)

	output, err := FindAuthorizedPrincipal(ctx, conn, d.Get("domain_name").(string), d.Id())

	return diags
}

func FindAuthorizedPrincipal(ctx context.Context, conn *opensearchservice.OpenSearchService, domainName string, id string) (*opensearchservice.AuthorizedPrincipal, error) {
    input := &opensearchservice.ListVpcEndpointAccessInput{
		DomainName: aws.String(domainName),
	}

	output, err := conn.ListVpcEndpointAccess(input)

	if err != nil {
		return nil, err
	}

	for _, principal := range output.AuthorizedPrincipalList {
		if id == "authorized-principal-" + *principal.Principal + "-" + *principal.PrincipalType + "-" + domainName {
			return principal, nil
		}
	}

	return nil, nil
}