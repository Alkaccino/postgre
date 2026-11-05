UPDATE assets
SET asset_name = @asset_name, category = @category, model = @model, manufacturer = @manufacturer, assigned_user = @assigned_user, asset_location = @asset_location, purchased_price = @purchased_price, purchased_at = @purchased_at, quantity = @quantity, asset_status = @asset_status, updated_at = @updated_at, created_by = @created_by, updated_by = @updated_by
WHERE id = @id