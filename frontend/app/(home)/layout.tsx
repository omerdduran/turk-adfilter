import type { ReactNode } from "react";
import { HomeLayout } from "fumadocs-ui/layouts/home";
import { baseOptions } from "@/app/layout.config";
import Footer from "@/components/footer";

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <div className="[&_#nd-nav>div]:bg-white/95 [&_#nd-nav>div]:dark:bg-[#191919]/95 [&_#nd-nav>div]:backdrop-blur-sm [&_#nd-nav>div]:border-b [&_#nd-nav>div]:border-gray-200 [&_#nd-nav>div]:dark:border-gray-800 [&_#nd-nav>div]:shadow-sm">
      <HomeLayout {...baseOptions}>
        {children}
        <Footer />
      </HomeLayout>
    </div>
  );
}
