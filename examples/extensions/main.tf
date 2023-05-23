terraform {
  required_providers {
    excel = {
      source = "deathfireofdoom.com/edu/excel"
    }
  }
}

provider "excel" {
}

data "excel_extensions" "edu" {}

output "excel_extensions" {
  value = data.excel_extensions.edu
}
