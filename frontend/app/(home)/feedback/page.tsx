import IssueForm from "@/components/issue-form";

export default function FeedbackPage() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-background to-muted/20 text-center">
      <div className="text-center mb-16 pt-16">
        <h1 className="text-4xl font-bold mb-4 bg-gradient-to-r from-red-600 to-red-700 dark:from-red-500 dark:to-red-400 bg-clip-text text-transparent">
          Geri Bildirimleriniz bizim için önemli!
        </h1>
        <p className="text-lg text-muted-foreground">
          Turk-AdFilter ile ilgili herhangi bir geri bildirim veya öneriniz varsa lütfen bize ulaşın.
        </p>
      </div>
      <IssueForm />
    </div>
  );
}
