package google_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleLoggingSink_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name": acctest.GetTestProjectFromEnv(),
		"sink_name":    "tf-test-sink-ds-" + acctest.RandString(t, 10),
		"bucket_name":  "tf-test-sink-ds-bucket-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:  func() { acctest.TestAccPreCheck(t) },
		Providers: acctest.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleLoggingSink_basic(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_logging_sink.basic",
						"google_logging_project_sink.basic",
						map[string]struct{}{
							"project":                {},
							"unique_writer_identity": {},
						},
					),
				),
			},
		},
	})
}

func testAccDataSourceGoogleLoggingSink_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_logging_project_sink" "basic" {
  name        = "%{sink_name}"
  project     = "%{project_name}"
  destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter      = "logName=\"projects/%{project_name}/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"

  unique_writer_identity = false
}

resource "google_storage_bucket" "log-bucket" {
  name     = "%{bucket_name}"
  location = "US"
}

data "google_logging_sink" "basic" {
  id = google_logging_project_sink.basic.id
}
`, context)
}
