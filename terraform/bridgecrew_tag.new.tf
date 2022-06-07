resource "bridgecrew_tag" "new" {
  name          = "james2"
  description   = "first terraform"
  isenabled     = false
  tagruleootbid = ""

  definition {
    tag_groups {
      name = "1653659061445_key"
      tags {
        value = {
          default = "name"
        }
        name = "team2"
      }
    }
  }

  repositories = ["00ca7905-d366-470e-9740-3a576fd9b82d", "02ecf59e-6cd6-4b14-9c91-816c46211bd2"]
}
