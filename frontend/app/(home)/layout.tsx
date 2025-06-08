import type { ReactNode } from "react";
import { HomeLayout } from "fumadocs-ui/layouts/home";
import { baseOptions } from "@/app/layout.config";
import Footer from "@/components/footer";

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <div className="[&_.fumadocs-nav]:bg-white/95 [&_.fumadocs-nav]:dark:bg-[#191919]/95 [&_.fumadocs-nav]:backdrop-blur-sm [&_.fumadocs-nav]:border-b [&_.fumadocs-nav]:border-gray-200 [&_.fumadocs-nav]:dark:border-gray-800 [&_.fumadocs-nav]:shadow-sm">
      <HomeLayout {...baseOptions}>
        {children}
        <Footer />
      </HomeLayout>
    </div>
  );
}
