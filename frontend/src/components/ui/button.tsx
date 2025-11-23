import * as React from "react";
import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";

import { cn } from "@/src/lib/utils";

const defaultButtonStyles =
  "mt-8 w-56 rounded-xl bg-[#6b5337] py-3 text-lg font-black text-white shadow-[0_6px_0_rgba(0,0,0,0.35)] hover:bg-[#7e6243] hover:-translate-y-0.5 hover:shadow-[0_7px_0_rgba(0,0,0,0.35)] active:translate-y-1 active:shadow-[0_2px_0_rgba(0,0,0,0.35)]";

const buttonVariants = cva(
  // ここは共通のレイアウト・状態系クラス（見た目ではなく挙動寄り）
  "inline-flex select-none items-center justify-center gap-2 whitespace-nowrap transition-all disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 shrink-0 [&_svg]:shrink-0 outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive",
  {
    variants: {
      variant: {
        default: defaultButtonStyles,
        primary: cn(
          defaultButtonStyles,
          "bg-button-primary hover:bg-button-primary/90"
        ),
        secondary: cn(
          defaultButtonStyles,
          "bg-button-secondary hover:bg-button-secondary/90"
        ),
        destructive: cn(
          defaultButtonStyles,
          "bg-destructive hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40 dark:bg-destructive/60"
        ),
        outline: cn(
          defaultButtonStyles,
          "border !bg-background !text-foreground shadow-xs hover:bg-accent hover:text-accent-foreground dark:bg-input/30 dark:border-input dark:hover:bg-input/50"
        ),
        ghost: cn(
          defaultButtonStyles,
          "hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50"
        ),
        link: cn(
          defaultButtonStyles,
          "text-primary underline-offset-4 hover:underline"
        ),
      },
      size: {
        // default サイズには余計な高さ・padding を入れないようにして、
        // 上の variant.default の py-3 を素直に効かせる
        default: "h-8 rounded-md gap-1.5 px-3 has-[>svg]:px-2.5 text-sm",
        sm: "h-8 rounded-md gap-1.5 px-3 has-[>svg]:px-2.5 text-sm",
        lg: "h-10 rounded-md px-6 has-[>svg]:px-4 text-base",
        icon: "size-9",
        "icon-sm": "size-8",
        "icon-lg": "size-10",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  }
);

function Button({
  className,
  variant,
  size,
  asChild = false,
  ...props
}: React.ComponentProps<"button"> &
  VariantProps<typeof buttonVariants> & {
    asChild?: boolean;
  }) {
  const Comp = asChild ? Slot : "button";

  return (
    <Comp
      data-slot="button"
      className={cn(
        buttonVariants({ variant, size, className }),
        "cursor-pointer"
      )}
      {...props}
    />
  );
}

export { Button, buttonVariants };
