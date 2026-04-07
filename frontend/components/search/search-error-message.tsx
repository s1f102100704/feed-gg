type SearchErrorMessageProps = {
  message: string;
  centered?: boolean;
};

export function SearchErrorMessage({
  message,
  centered = false,
}: SearchErrorMessageProps) {
  return (
    <p className={centered ? "text-center text-sm text-red-300" : "text-sm text-red-300"}>
      {message}
    </p>
  );
}
