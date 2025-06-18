export type UploadedFileProps = {
  id?: number;
  publicId?: string;
  secureUrl?: string;
  format?: string;
  resourceType?: string;
  size?: number;
  width?: number;
  height?: number;
  originalFilename?: string;
  createdAt?: string;
  updatedAt?: string;
  isUploading?: boolean;
};

export type UploadedFileResponseProps = {
  success: string;
};
