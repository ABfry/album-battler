import type { Meta, StoryObj } from "@storybook/react";
import { Label } from "./label";
import { Input } from "./input";

const meta = {
  title: "UI/Label",
  component: Label,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
} satisfies Meta<typeof Label>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    children: "Label Text",
  },
};

export const WithInput: Story = {
  render: () => (
    <div className="w-[350px] space-y-2">
      <Label htmlFor="email">メールアドレス</Label>
      <Input id="email" type="email" placeholder="email@example.com" />
    </div>
  ),
};

export const Required: Story = {
  render: () => (
    <div className="w-[350px] space-y-2">
      <Label htmlFor="required-input">
        必須項目
        <span className="text-destructive ml-1">*</span>
      </Label>
      <Input id="required-input" type="text" placeholder="入力してください" />
    </div>
  ),
};

export const WithHelpText: Story = {
  render: () => (
    <div className="w-[350px] space-y-2">
      <Label htmlFor="username">ユーザー名</Label>
      <Input id="username" type="text" placeholder="ユーザー名を入力" />
      <p className="text-muted-foreground text-sm">
        半角英数字とアンダースコアのみ使用できます
      </p>
    </div>
  ),
};

export const MultipleFields: Story = {
  render: () => (
    <div className="w-[350px] space-y-4">
      <div className="space-y-2">
        <Label htmlFor="field1">名前</Label>
        <Input id="field1" type="text" placeholder="山田太郎" />
      </div>
      <div className="space-y-2">
        <Label htmlFor="field2">メールアドレス</Label>
        <Input id="field2" type="email" placeholder="email@example.com" />
      </div>
      <div className="space-y-2">
        <Label htmlFor="field3">電話番号</Label>
        <Input id="field3" type="tel" placeholder="090-1234-5678" />
      </div>
    </div>
  ),
};

export const Disabled: Story = {
  render: () => (
    <div className="w-[350px] space-y-2">
      <Label htmlFor="disabled-input" className="opacity-50">
        無効化された入力
      </Label>
      <Input
        id="disabled-input"
        type="text"
        placeholder="入力できません"
        disabled
      />
    </div>
  ),
};

export const WithDescription: Story = {
  render: () => (
    <div className="w-[350px] space-y-2">
      <div className="flex items-center justify-between">
        <Label htmlFor="password">パスワード</Label>
        <span className="text-muted-foreground text-sm">8文字以上</span>
      </div>
      <Input id="password" type="password" placeholder="パスワードを入力" />
    </div>
  ),
};
