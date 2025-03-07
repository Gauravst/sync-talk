import ChatPreviewSection from "@/components/home/ChatPreviewSection";
import CTASection from "@/components/home/CTASection";
import FeaturesSection from "@/components/home/FeaturesSection";
import Footer from "@/components/home/Footer";
import Header from "@/components/home/Header";
import HeroSection from "@/components/home/HeroSection";

function HomePage() {
  return (
    <div className="flex flex-col w-full justify-center items-center min-h-screen bg-background">
      {/* Header */}
      <Header />

      {/* Hero Section */}
      <HeroSection />

      {/* Chat Rooms Preview Section */}
      <ChatPreviewSection />

      {/* Features Section */}
      <FeaturesSection />

      {/* CTA Section */}
      <CTASection />

      {/* Footer */}
      <Footer />
    </div>
  );
}

export default HomePage;
