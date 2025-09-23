-- name: CreateSupplier :execrows
INSERT INTO dim_supplier_v1
    (
    supplier_id,
    supplier_code,
    tenant_id,
    legal_name,
    dba_name,
    country,
    region,
    address_line1,
    address_line2,
    city,
    state,
    postal_code,
    contact_email,
    contact_phone,
    preferred_currency,
    incoterms,
    lead_time_days_avg,
    lead_time_days_p95,
    on_time_delivery_rate,
    defect_rate_ppm,
    capacity_units_per_week,
    risk_score,
    financial_risk_tier,
    certifications,
    compliance_flags,
    approved_status,
    contracts,
    terms_version,
    lat,
    lon,
    data_source,
    source_timestamp,
    ingestion_timestamp,
    schema_version
    )
VALUES
    (
        ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
        ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
        ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
        ?, ?, ?, ?
);

-- name: CreatePart :execrows
INSERT INTO dim_part_v1
    (
    part_id,
    tenant_id,
    part_number,
    description,
    category,
    lifecycle_status,
    uom,
    spec_hash,
    bom_compatibility,
    default_supplier_id,
    qualified_supplier_ids,
    unit_cost,
    moq,
    lead_time_days_avg,
    lead_time_days_p95,
    quality_grade,
    compliance_flags,
    hazard_class,
    last_price_change,
    data_source,
    source_timestamp,
    ingestion_timestamp,
    schema_version
    )
VALUES
    (
        ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
        ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
        ?, ?, ?
);

-- name: DeleteSupplier :exec
DELETE FROM dim_supplier_v1 WHERE supplier_id = ?;

-- name: DeletePart :exec
DELETE FROM dim_part_v1 WHERE part_id = ?;
