type SelectImagePreviewProps = {
  url: string;
  open: boolean;
  close: () => void;
};

export const SelectImagePreview = ({ url, open }: SelectImagePreviewProps) => {
  if (!open || !url) return null;
  return (
    <div className="absolute z-50 h-[100%] inset-0 w-full bg-background flex items-center justify-center ">
      <img
        src={url}
        className={`object-contain w-auto max-h-[80%] rounded-2xl`}
      />
    </div>
  );
};
