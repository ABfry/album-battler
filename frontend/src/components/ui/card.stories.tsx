import type { Meta, StoryObj } from "@storybook/react";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
  CardAction,
} from "./card";
import { Button } from "./button";

const meta = {
  title: "UI/Card",
  component: Card,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
} satisfies Meta<typeof Card>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => (
    <Card className="w-[350px]">
      <CardHeader>
        <CardTitle>Card Title</CardTitle>
        <CardDescription>Card Description</CardDescription>
      </CardHeader>
      <CardContent>
        <p>Card content goes here.</p>
      </CardContent>
    </Card>
  ),
};

export const WithFooter: Story = {
  render: () => (
    <Card className="w-[350px]">
      <CardHeader>
        <CardTitle>Card with Footer</CardTitle>
        <CardDescription>This card has a footer section</CardDescription>
      </CardHeader>
      <CardContent>
        <p>Card content goes here.</p>
      </CardContent>
      <CardFooter>
        <Button variant="primary" className="w-full">
          Action
        </Button>
      </CardFooter>
    </Card>
  ),
};

export const WithAction: Story = {
  render: () => (
    <Card className="w-[350px]">
      <CardHeader>
        <CardTitle>Card with Action</CardTitle>
        <CardDescription>
          This card has an action button in header
        </CardDescription>
        <CardAction>
          <Button variant="ghost" size="sm">
            Edit
          </Button>
        </CardAction>
      </CardHeader>
      <CardContent>
        <p>Card content goes here.</p>
      </CardContent>
    </Card>
  ),
};

export const FormCard: Story = {
  render: () => (
    <Card className="w-[350px]">
      <CardHeader>
        <CardTitle>Create Account</CardTitle>
        <CardDescription>
          Enter your information to create an account
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <label className="text-sm font-medium">Name</label>
          <input
            type="text"
            placeholder="Enter your name"
            className="w-full rounded-md border px-3 py-2"
          />
        </div>
        <div className="space-y-2">
          <label className="text-sm font-medium">Email</label>
          <input
            type="email"
            placeholder="Enter your email"
            className="w-full rounded-md border px-3 py-2"
          />
        </div>
      </CardContent>
      <CardFooter className="gap-2">
        <Button variant="outline" className="w-full">
          Cancel
        </Button>
        <Button variant="primary" className="w-full">
          Create
        </Button>
      </CardFooter>
    </Card>
  ),
};

export const InfoCard: Story = {
  render: () => (
    <Card className="w-[350px]">
      <CardHeader>
        <CardTitle>アルバムバトラー</CardTitle>
        <CardDescription>部屋番号: 12345</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          <div className="flex items-center gap-3">
            <div>
              <p className="font-bold">岩崎</p>
              <p className="text-muted-foreground text-sm">プレイヤー1</p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <div>
              <p className="font-bold">井上</p>
              <p className="text-muted-foreground text-sm">プレイヤー2</p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <div>
              <p className="font-bold">シバタ</p>
              <p className="text-muted-foreground text-sm">プレイヤー3</p>
            </div>
          </div>
        </div>
      </CardContent>
      <CardFooter>
        <Button variant="secondary" className="w-full">
          ゲーム開始
        </Button>
      </CardFooter>
    </Card>
  ),
};
