import { UploadedFileProps } from "@/types/fileTypes";

type UploadImagePreviewProps = {
  file: UploadedFileProps;
  isUploading: boolean;
  progress: number;
};

export const ImagePreview = ({
  file,
  isUploading,
  progress,
}: UploadImagePreviewProps) => {
  console.log("file--", file);
  return (
    <div className="relative w-48 h-48 rounded overflow-hidden cursor-pointer">
      <img
        src={file.secureUrl}
        className={`object-cover w-full h-full ${isUploading && "opacity-70"}`}
      />

      {/* Loading Overlay */}
      {isUploading && (
        <div className="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50">
          <div className="w-32 h-2 bg-gray-300 rounded">
            <div
              className="h-2 bg-blue-500 rounded"
              style={{ width: `${progress}%` }}
            ></div>
          </div>
        </div>
      )}
    </div>
  );
};
