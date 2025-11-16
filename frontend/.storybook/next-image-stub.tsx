import * as React from "react";

// Storybook用のNext.js Imageコンポーネント
const NextImage = ({
  src,
  alt,
  width,
  height,
  className,
  ...props
}: React.ComponentProps<"img"> & {
  src: string;
  alt: string;
  width?: number | string;
  height?: number | string;
  priority?: boolean;
  loading?: "lazy" | "eager";
}) => {
  return (
    // eslint-disable-next-line @next/next/no-img-element
    <img
      src={src}
      alt={alt}
      width={typeof width === "number" ? width : undefined}
      height={typeof height === "number" ? height : undefined}
      className={className}
      {...props}
    />
  );
};

export default NextImage;
