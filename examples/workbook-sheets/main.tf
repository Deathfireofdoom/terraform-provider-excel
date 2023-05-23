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
  value = resource.excel_workbook.edu
}

resource "excel_workbook" "edu" {
    file_name = "edu"
    folder_path = "workbooks"
    extension = "xlsx"
}

resource "excel_sheet" "edu" {
  workbook_id = excel_workbook.edu.id
  name = "edu_test"
}