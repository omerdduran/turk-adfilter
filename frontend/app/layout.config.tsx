import { GithubInfo } from "fumadocs-ui/components/github-info";
import type { BaseLayoutProps } from "fumadocs-ui/layouts/shared";

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
          style={{ display: "inline", verticalAlign: "middle", marginRight: 8 }}
        />
        <span style={{ verticalAlign: "middle" }}>Turk-AdFilter</span>
      </>
    ),
  },
  links: [
    {
      text: "Hakkımızda",
      url: "/about",
      active: "nested-url",
      icon: "❓ ",
    },
    {
      text: "Dokümantasyon",
      url: "/docs",
      active: "nested-url",
      icon: "📚 ",
    },
    {
      text: "Sürüm Notları",
      url: "/releases",
      active: "nested-url",
      icon: "📋 ",
    },
    {
      text: "Filtre Yedekleri",
      url: "/backups",
      active: "nested-url",
      icon: "💾 ",
    },
    {
      text: "Geri Bildirim",
      url: "/feedback",
      active: "nested-url",
      icon: "💬 ",
    },
    {
      type: "icon",
      url: "https://github.com/omerdduran/turk-adfilter",
      text: "Github",
      icon: (
        <svg role="img" viewBox="0 0 24 24" fill="currentColor">
          <path d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12" />
        </svg>
      ),
      external: true,
    },
    {
      type: "icon",
      url: "https://matrix.to/#/#reklamsiz-turkiye:matrix.org",
      text: "Matrix",
      icon: (
        <svg
          version="1.1"
          id="Layer_1"
          xmlns="http://www.w3.org/2000/svg"
          x="0px"
          y="0px"
          viewBox="0 0 520 520"
          fill="currentColor"
        >
          <path d="M13.7,11.9v496.2h35.7V520H0V0h49.4v11.9H13.7z" />
          <path d="M166.3,169.2v25.1h0.7c6.7-9.6,14.8-17,24.2-22.2c9.4-5.3,20.3-7.9,32.5-7.9c11.7,0,22.4,2.3,32.1,6.8  c9.7,4.5,17,12.6,22.1,24c5.5-8.1,13-15.3,22.4-21.5c9.4-6.2,20.6-9.3,33.5-9.3c9.8,0,18.9,1.2,27.3,3.6c8.4,2.4,15.5,6.2,21.5,11.5  c6,5.3,10.6,12.1,14,20.6c3.3,8.5,5,18.7,5,30.7v124.1h-50.9V249.6c0-6.2-0.2-12.1-0.7-17.6c-0.5-5.5-1.8-10.3-3.9-14.3  c-2.2-4.1-5.3-7.3-9.5-9.7c-4.2-2.4-9.9-3.6-17-3.6c-7.2,0-13,1.4-17.4,4.1c-4.4,2.8-7.9,6.3-10.4,10.8c-2.5,4.4-4.2,9.4-5,15.1  c-0.8,5.6-1.3,11.3-1.3,17v103.3h-50.9v-104c0-5.5-0.1-10.9-0.4-16.3c-0.2-5.4-1.3-10.3-3.1-14.9c-1.8-4.5-4.8-8.2-9-10.9  c-4.2-2.7-10.3-4.1-18.5-4.1c-2.4,0-5.6,0.5-9.5,1.6c-3.9,1.1-7.8,3.1-11.5,6.1c-3.7,3-6.9,7.3-9.5,12.9c-2.6,5.6-3.9,13-3.9,22.1  v107.6h-50.9V169.2H166.3z" />
          <path d="M506.3,508.1V11.9h-35.7V0H520v520h-49.4v-11.9H506.3z" />
        </svg>
      ),
      external: true,
    },
  ],
};
