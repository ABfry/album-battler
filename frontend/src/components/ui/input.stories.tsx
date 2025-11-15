import type { Meta, StoryObj } from "@storybook/react";
import { Input } from "./input";
import { Label } from "./label";

const meta = {
  title: "UI/Input",
  component: Input,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
  argTypes: {
    type: {
      control: "select",
      options: ["text", "email", "password", "number", "tel", "url", "search"],
    },
    disabled: {
      control: "boolean",
    },
    placeholder: {
      control: "text",
    },
  },
} satisfies Meta<typeof Input>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    placeholder: "Enter text...",
  },
};

export const Text: Story = {
  args: {
    type: "text",
    placeholder: "テキストを入力",
  },
};

export const Email: Story = {
  args: {
    type: "email",
    placeholder: "email@example.com",
  },
};

export const Password: Story = {
  args: {
    type: "password",
    placeholder: "パスワードを入力",
  },
};

export const Number: Story = {
  args: {
    type: "number",
    placeholder: "番号を入力",
  },
};

export const Disabled: Story = {
  args: {
    placeholder: "Disabled input",
    disabled: true,
    value: "This input is disabled",
  },
};

export const WithLabel: Story = {
  render: () => (
    <div className="w-[350px] space-y-2">
      <Label htmlFor="input-with-label">メールアドレス</Label>
      <Input
        id="input-with-label"
        type="email"
        placeholder="email@example.com"
      />
    </div>
  ),
};

export const FormExample: Story = {
  render: () => (
    <div className="w-[350px] space-y-4">
      <div className="space-y-2">
        <Label htmlFor="username">ユーザー名</Label>
        <Input id="username" type="text" placeholder="ユーザー名を入力" />
      </div>
      <div className="space-y-2">
        <Label htmlFor="email">メールアドレス</Label>
        <Input id="email" type="email" placeholder="email@example.com" />
      </div>
      <div className="space-y-2">
        <Label htmlFor="password">パスワード</Label>
        <Input id="password" type="password" placeholder="パスワードを入力" />
      </div>
      <div className="space-y-2">
        <Label htmlFor="room-number">部屋番号</Label>
        <Input id="room-number" type="number" placeholder="12345" />
      </div>
    </div>
  ),
};

export const Search: Story = {
  args: {
    type: "search",
    placeholder: "検索...",
  },
};

export const WithValue: Story = {
  args: {
    type: "text",
    value: "Pre-filled value",
    placeholder: "Enter text...",
  },
};
