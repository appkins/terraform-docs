/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package terraform

import "github.com/hashicorp/hcl/v2"

// Position represents position of Terraform item (input, output, provider, etc) in a file.
type Position *hcl.Range
