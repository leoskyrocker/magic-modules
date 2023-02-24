package google_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeAddress_networkTier(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.TestAccPreCheck(t) },
		Providers:    acctest.TestAccProviders,
		CheckDestroy: testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeAddress_networkTier(acctest.RandString(t, 10)),
			},
			{
				ResourceName:      "google_compute_address.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeAddress_internal(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.TestAccPreCheck(t) },
		Providers:    acctest.TestAccProvidersrovidersroviders,
		CheckDestroy: testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeAddress_internal(acctest.RandString(t, 10)),
			},
			{
				ResourceName:      "google_compute_address.internal",
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				ResourceName:      "google_compute_address.internal_with_subnet",
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				ResourceName:      "google_compute_address.internal_with_subnet_and_address",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeAddress_internal(i string) string {
	return fmt.Sprintf(`
resource "google_compute_address" "internal" {
  name         = "tf-test-address-internal-%s"
  address_type = "INTERNAL"
  region       = "us-east1"
}

resource "google_compute_network" "default" {
  name = "tf-test-network-test-%s"
}

resource "google_compute_subnetwork" "foo" {
  name          = "subnetwork-test-%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-east1"
  network       = google_compute_network.default.self_link
}

resource "google_compute_address" "internal_with_subnet" {
  name         = "tf-test-address-internal-with-subnet-%s"
  subnetwork   = google_compute_subnetwork.foo.self_link
  address_type = "INTERNAL"
  region       = "us-east1"
}

// We can't test the address alone, because we don't know what IP range the
// default subnetwork uses.
resource "google_compute_address" "internal_with_subnet_and_address" {
  name         = "tf-test-address-internal-with-subnet-and-address-%s"
  subnetwork   = google_compute_subnetwork.foo.self_link
  address_type = "INTERNAL"
  address      = "10.0.42.42"
  region       = "us-east1"
}
`,
		i, // google_compute_address.internal name
		i, // google_compute_network.default name
		i, // google_compute_subnetwork.foo name
		i, // google_compute_address.internal_with_subnet_name
		i, // google_compute_address.internal_with_subnet_and_address name
	)
}

func testAccComputeAddress_networkTier(i string) string {
	return fmt.Sprintf(`
resource "google_compute_address" "foobar" {
  name         = "tf-test-address-%s"
  network_tier = "STANDARD"
}
`, i)
}

func TestAccProviderBasePath_setBasePath(t *testing.T) {
	t.Parallel()

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProviderBasePath_setBasePath("https://www.googleapis.com/compute/beta/", RandString(t, 10)),
			},
			{
				ResourceName:      "google_compute_address.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Move this test here because the testAccCheckComputeAddressDestroyProducer is used to check destroy.
func TestAccProviderBasePath_setInvalidBasePath(t *testing.T) {
	t.Parallel()

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccProviderBasePath_setBasePath("https://www.example.com/compute/beta/", RandString(t, 10)),
				ExpectError: regexp.MustCompile("got HTTP response code 404 with body"),
			},
		},
	})
}

func TestAccProviderMeta_setModuleName(t *testing.T) {
	t.Parallel()

	moduleName := "my-module"
	VcrTest(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProviderMeta_setModuleName(moduleName, RandString(t, 10)),
			},
			{
				ResourceName:      "google_compute_address.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProviderBasePath_setBasePath(endpoint, name string) string {
	return fmt.Sprintf(`
provider "google" {
  alias                   = "compute_custom_endpoint"
  compute_custom_endpoint = "%s"
}

resource "google_compute_address" "default" {
  provider = google.compute_custom_endpoint
  name     = "tf-test-address-%s"
}`, endpoint, name)
}

func testAccProviderMeta_setModuleName(key, name string) string {
	return fmt.Sprintf(`
terraform {
  provider_meta "google" {
    module_name = "%s"
  }
}

resource "google_compute_address" "default" {
	name = "tf-test-address-%s"
}`, key, name)
}
