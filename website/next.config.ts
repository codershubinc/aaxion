import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async redirects() {
    return [
      {
        source: '/install.sh',
        destination: 'https://raw.githubusercontent.com/codershubinc/aaxion/main/install.sh',
        permanent: false,
      },
    ]
  },
};
module.exports = {
  allowedDevOrigins: ['192.168.1.107'],
}

export default nextConfig;
