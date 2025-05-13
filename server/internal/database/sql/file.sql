-- name: UploadFile
INSERT INTO
  files (
    publicId,
    secureURL,
    format,
    resourceType,
    size,
    width,
    height,
    originalFilename
  )
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
  id,
  publicId,
  secureUrl,
  format,
  resourceType,
  size,
  width,
  height,
  originalFilename,
  createdAt
