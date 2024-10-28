schema "public" {
  comment = "standard public schema"
}

table "function_histories" {
  schema = schema.public

  column "id" {
    null = false
    type = bigserial
  }
  column "associate_with" {
    null = false
    type = varchar
  }
  column "called_by_function" {
    null = false
    type = varchar
  }
  column "line" {
    null = false
    type = bigint
  }
  column "file_location" {
    null = false
    type = text
  }

  column "created_at" {
    null    = false
    type    = timestamp(3)
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

}

table "discord_toggle_histories" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "action_by" {
     null = false
     type = bigserial
  }
  column "payload" {
     null = false
     type = json
  }

  column "created_at" {
    null    = false
    type    = timestamp(3)
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }
  index "ix_discord_toggle_histories_action_by" {
    columns = [column.action_by]
  }
}

table "user_otps" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigserial
  }
  column "otp" {
    null = false
    type = varchar
  }
  column "expired_at" {
      null    = false
      type    = timestamp(3)
      default = sql("CURRENT_TIMESTAMP + INTERVAL '1 hour'")
  }
  column "created_at" {
    null    = false
    type    = timestamp(3)
    default = sql("CURRENT_TIMESTAMP")
  }

  foreign_key "user_id_fk" {
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete = CASCADE
    on_update = NO_ACTION
  }
}

table "users" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "email" {
     null = true
     type = varchar
  }
  column "phone_number" {
     null = true
     type = varchar
  }
  column "password_hash" {
    null = false
    type = varchar
  }
  column "admin" {
     null = false
     type = bool
     default = false
  }
  column "verify_by" {
     null = false
     type = int
     default = 0
  }

  column "created_at" {
    null    = false
    type    = timestamp(3)
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = false
    type    = timestamp(3)
    default = sql("CURRENT_TIMESTAMP")
  }

  column "deleted_at" {
    null    = true
    type    = timestamp(3)
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_users_email" {
    columns = [column.email]
  }

  index "idx_users_phone" {
     columns = [column.phone_number]
  }

  index "unique_users_email" {
      columns = [column.email]
      where = "deleted_at IS NULL AND email != ''"
      unique = true
  }

  index "unique_users_phone" {
      columns = [column.phone_number]
      where = "deleted_at IS NULL AND phone_number != ''"
      unique = true
  }

}
