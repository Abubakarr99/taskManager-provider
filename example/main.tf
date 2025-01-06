terraform {
  required_providers {
    taskmanager = {
      version = "0.0.1"
      source  = "dantata.com/aboudev/taskmanager"
    }
  }
}

provider "taskmanager" {
  host = "127.0.0.1:6742"
}

resource "taskmanager_task" "one" {
  title    = "send email tomorrow"
  urgency  = "VERY_HIGH"
  due_date = "2025-01-06"
}

