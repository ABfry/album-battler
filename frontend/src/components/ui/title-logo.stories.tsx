import type { Meta, StoryObj } from "@storybook/react";
import { TitleLogo } from "./title-logo";

const meta = {
  title: "UI/TitleLogo",
  component: TitleLogo,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
  argTypes: {
    alt: {
      control: "text",
      description: "代替テキスト",
    },
    className: {
      control: "text",
      description: "CSSクラス（w-*, h-*でサイズ指定）",
    },
    loading: {
      control: "select",
      options: ["lazy", "eager"],
      description: "画像の読み込み方法",
    },
  },
} satisfies Meta<typeof TitleLogo>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    className: "w-64",
  },
};

export const CustomHeight: Story = {
  args: {
    className: "h-32",
  },
};

export const EagerLoading: Story = {
  args: {
    className: "w-64",
    loading: "eager",
  },
};

export const FullWidth: Story = {
  args: {
    className: "w-full",
  },
  render: (args) => (
    <div className="w-96 border border-gray-300">
      <p className="mb-2 text-sm text-gray-600">
        親要素の幅に合わせる (className=&quot;w-full&quot;)
      </p>
      <TitleLogo {...args} />
    </div>
  ),
};

export const FullHeight: Story = {
  args: {
    className: "h-full",
  },
  render: (args) => (
    <>
      <p className="mb-2 text-sm text-gray-600">
        親要素の高さに合わせる (className=&quot;h-full&quot;)
      </p>
      <div className="h-64 border border-gray-300">
        <TitleLogo {...args} />
      </div>
    </>
  ),
};
