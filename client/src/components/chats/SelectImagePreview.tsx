type SelectImagePreviewProps = {
  url: string;
  open: boolean;
  close: () => void;
};

export const SelectImagePreview = ({ url, open }: SelectImagePreviewProps) => {
  if (!open || !url) return null;
  return (
    <div className="absolute z-50 h-[calc(100%-60px)] inset-0 mt-[60px] w-full bg-background flex items-center justify-center ">
      <img
        src={url}
        className={`object-contain w-auto max-h-[80%] rounded-2xl`}
      />
    </div>
  );
};
