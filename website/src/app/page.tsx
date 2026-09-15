import Link from 'next/link';
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { InstallCommand } from '@/components/install-command';

export default function Home() {
  return (
    <div className="min-h-screen bg-black text-white selection:bg-blue-500/30 relative overflow-hidden">

      {/* Background Glow Effects */}
      <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] rounded-full bg-blue-600/20 blur-[120px] pointer-events-none" />
      <div className="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] rounded-full bg-indigo-600/20 blur-[120px] pointer-events-none" />

      {/* Navbar */}
      <nav className="border-b border-white/5 sticky top-0 z-50 bg-black/50 backdrop-blur-xl">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16 items-center">
            <div className="flex items-center space-x-2">
              <span className="text-2xl font-bold bg-linear-to-r from-blue-400 to-indigo-500 bg-clip-text text-transparent">Aaxion</span>
            </div>
            <div className="hidden md:flex space-x-8 items-center">
              <Link href="#features" className="text-sm font-medium text-gray-400 hover:text-white transition-colors">Features</Link>
              <Link href="#installation" className="text-sm font-medium text-gray-400 hover:text-white transition-colors">Installation</Link>
              <Link href="/docs" className="text-sm font-medium text-gray-400 hover:text-white transition-colors">Docs</Link>
              <a href="https://github.com/codershubinc/aaxion" target="_blank" rel="noreferrer">
                <Button variant="outline" size="sm" className="bg-white/5 border-white/10 hover:bg-white/10 hover:text-white transition-all">GitHub &rarr;</Button>
              </a>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-5xl mx-auto px-4 py-24 relative z-10">

        {/* Hero Section */}
        <div className="text-center mb-32 flex flex-col items-center pt-10">
          <Badge variant="outline" className="mb-6 bg-blue-500/10 text-blue-400 border-blue-500/20 px-4 py-1.5 text-sm rounded-full backdrop-blur-md">
            v0.0.1-beta &quot;Photon&quot; is live
          </Badge>
          <h1 className="text-6xl md:text-8xl font-extrabold mb-8 tracking-tighter leading-tight">
            Zero-Buffer <br />
            <span className="bg-linear-to-r from-blue-400 via-indigo-400 to-purple-400 bg-clip-text text-transparent">File Streaming.</span>
          </h1>
          <p className="text-xl md:text-2xl text-gray-400 mb-12 max-w-3xl mx-auto font-light leading-relaxed">
            Repurpose your old hardware into blazing-fast storage nodes. 10GB transfers using just *~32KB of RAM.
          </p>
          <div className="flex flex-col sm:flex-row justify-center gap-4 w-full sm:w-auto mb-10">
            <Link href="#installation" className="w-full sm:w-auto">
              <Button size="lg" className="w-full sm:w-auto font-semibold bg-white text-black hover:bg-gray-200 shadow-[0_0_40px_-10px_rgba(255,255,255,0.3)] h-14 px-8 text-base rounded-xl transition-all">
                Quick Install
              </Button>
            </Link>
            <Link href="/docs" className="w-full sm:w-auto">
              <Button size="lg" variant="outline" className="w-full sm:w-auto font-semibold bg-white/5 border-white/10 hover:bg-white/10 hover:text-white h-14 px-8 text-base rounded-xl backdrop-blur-md transition-all">
                Read the Docs
              </Button>
            </Link>
          </div>

          <InstallCommand variant="hero" />
        </div>

        {/* Features Section */}
        <section id="features" className="scroll-mt-32 mb-32">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-5xl font-bold tracking-tight mb-4">Engineered for Efficiency</h2>
            <p className="text-gray-400 text-lg">Every byte is streamed perfectly without wasting your memory.</p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <Card className="bg-white/5 border-white/10 backdrop-blur-lg hover:bg-white/[0.07] transition-colors">
              <CardHeader>
                <CardTitle className="text-2xl text-white">Direct I/O</CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription className="text-base text-gray-400 leading-relaxed">
                  Uploads and downloads stream directly to disk/network. Bypass memory bottlenecks completely.
                </CardDescription>
              </CardContent>
            </Card>

            <Card className="bg-white/5 border-white/10 backdrop-blur-lg hover:bg-white/[0.07] transition-colors">
              <CardHeader>
                <CardTitle className="text-2xl text-white">Resumable Chunks</CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription className="text-base text-gray-400 leading-relaxed">
                  Supports chunked uploading to bypass network limits (like Cloudflare Tunnels) and effortlessly resume interrupted transfers.
                </CardDescription>
              </CardContent>
            </Card>

            <Card className="bg-white/5 border-white/10 backdrop-blur-lg hover:bg-white/[0.07] transition-colors">
              <CardHeader>
                <CardTitle className="text-2xl text-white">One-Time Links</CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription className="text-base text-gray-400 leading-relaxed">
                  Generate hyper-secure, one-time temporary download links for external file sharing with zero risk.
                </CardDescription>
              </CardContent>
            </Card>

            <Card className="bg-white/5 border-white/10 backdrop-blur-lg hover:bg-white/[0.07] transition-colors relative overflow-hidden group">
              <div className="absolute inset-0 bg-linear-to-br from-blue-500/10 to-purple-500/10 opacity-0 group-hover:opacity-100 transition-opacity" />
              <CardHeader className="relative z-10">
                <CardTitle className="text-2xl text-white flex items-center gap-2">
                  <span className="relative flex h-3 w-3 mr-1">
                    <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                    <span className="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
                  </span>
                  ~10MB Idle RAM
                </CardTitle>
              </CardHeader>
              <CardContent className="relative z-10">
                <CardDescription className="text-base text-gray-400 leading-relaxed">
                  Consumes virtually zero resources when idle. It breathes life back into ancient, dusty hardware.
                </CardDescription>
              </CardContent>
            </Card>
          </div>
        </section>

        {/* Installation Section */}
        <section id="installation" className="scroll-mt-32 mb-32 relative">
          <div className="absolute inset-0 bg-blue-500/5 blur-[100px] pointer-events-none rounded-full" />
          <div className="bg-white/5 border border-white/10 rounded-3xl p-8 md:p-12 backdrop-blur-xl relative z-10">
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight mb-4 text-white">Install in Seconds</h2>
            <p className="mb-8 text-gray-400 text-lg max-w-2xl">Aaxion is distributed as a single binary. Run this command on your server to instantly install the latest release.</p>

            <InstallCommand variant="section" />

            <div className="mt-8 grid grid-cols-1 md:grid-cols-2 gap-8 pt-8 border-t border-white/5">
              <div>
                <h3 className="text-white font-semibold text-xl mb-3 flex items-center gap-3">
                  <span className="flex items-center justify-center w-8 h-8 rounded-full bg-blue-500/20 text-blue-400 text-sm border border-blue-500/30">1</span>
                  Create Admin
                </h3>
                <code className="block bg-black/40 border border-white/5 p-3 rounded-xl text-sm font-mono text-gray-300">
                  aax create-admin <span className="text-emerald-400">user:pass</span>
                </code>
              </div>
              <div>
                <h3 className="text-white font-semibold text-xl mb-3 flex items-center gap-3">
                  <span className="flex items-center justify-center w-8 h-8 rounded-full bg-indigo-500/20 text-indigo-400 text-sm border border-indigo-500/30">2</span>
                  Boot Server
                </h3>
                <code className="block bg-black/40 border border-white/5 p-3 rounded-xl text-sm font-mono text-gray-300">
                  aax serve
                </code>
              </div>
            </div>
          </div>
        </section>

      </main>

      {/* Footer */}
      <footer className="border-t border-white/5 py-12 text-center relative z-10 bg-black">
        <p className="text-gray-500 text-sm font-medium">&copy; 2026 CodersHub Inc. Open Source under the AGPLv3 License.</p>
      </footer>
    </div>
  );
}
