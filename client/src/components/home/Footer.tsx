const Footer = () => {
  return (
    <footer className="border-t bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container py-8">
        <div className="grid grid-cols-2 gap-8 md:grid-cols-4">
          <div className="space-y-4">
            <div className="flex items-center gap-2 font-bold text-xl">
              <img className="w-7 h-7" src="/icon.png" alt="logo" />
              <span>Sync Talk</span>
            </div>
            <p className="text-sm text-muted-foreground">
              Real-time chat powered by Go and WebSockets.
            </p>
          </div>
          <div className="space-y-3">
            <h4 className="font-medium">Product</h4>
            <ul className="space-y-2 text-sm text-muted-foreground">
              <li>Features</li>
              <li>Pricing</li>
              <li>Integrations</li>
              <li>Documentation</li>
            </ul>
          </div>
          <div className="space-y-3">
            <h4 className="font-medium">Company</h4>
            <ul className="space-y-2 text-sm text-muted-foreground">
              <li>About</li>
              <li>Blog</li>
              <li>Careers</li>
              <li>Contact</li>
            </ul>
          </div>
          <div className="space-y-3">
            <h4 className="font-medium">Legal</h4>
            <ul className="space-y-2 text-sm text-muted-foreground">
              <li>Terms</li>
              <li>Privacy</li>
              <li>Cookies</li>
              <li>Licenses</li>
            </ul>
          </div>
        </div>
        <div className="mt-8 border-t pt-8 text-center text-sm text-muted-foreground">
          © {new Date().getFullYear()} Sync Talk. All rights reserved.
        </div>
      </div>
    </footer>
  );
};

export default Footer;
