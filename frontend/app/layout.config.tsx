import type { BaseLayoutProps } from 'fumadocs-ui/layouts/shared';

/**
 * Shared layout configurations
 *
 * you can customise layouts individually from:
 * Home Layout: app/(home)/layout.tsx
 * Docs Layout: app/docs/layout.tsx
 */
export const baseOptions: BaseLayoutProps = {
  nav: {
    title: (
      <>
        <img
          src="/assets/logo.png"
          alt="Turk-AdFilter Logo"
          width={32}
          height={32}
          style={{ display: 'inline', verticalAlign: 'middle', marginRight: 8 }}
        />
        <span style={{ verticalAlign: 'middle' }}>Turk-AdFilter</span>
      </>
    ),
  },
  links: [
    {
      text: 'Documentation',
      url: '/docs',
      active: 'nested-url',
    },
  ],
};
