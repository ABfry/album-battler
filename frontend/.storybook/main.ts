import type { StorybookConfig } from "@storybook/react-vite";
import { join } from "path";

const config: StorybookConfig = {
  stories: ["../src/**/*.mdx", "../src/**/*.stories.@(js|jsx|mjs|ts|tsx)"],
  addons: ["@storybook/addon-essentials"],
  framework: {
    name: "@storybook/react-vite",
    options: {},
  },
  staticDirs: ["../public"],
  core: {
    disableTelemetry: true,
  },
  viteFinal: async (config) => {
    // Mock Next.js Image for Storybook
    config.resolve = config.resolve || {};
    config.resolve.alias = {
      ...config.resolve.alias,
      "next/image": join(__dirname, "next-image-stub.tsx"),
    };

    // Define process.env for Next.js compatibility
    config.define = {
      ...config.define,
      "process.env": {},
    };

    return config;
  },
};

export default config;
