const DISPATCH_URL =
  "https://api.github.com/repos/BelgianNoise/colruyt-products-scraper/actions/workflows/daily.yaml/dispatches";

export default {
  async scheduled(event, env, ctx) {
    const res = await fetch(DISPATCH_URL, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${env.GH_TOKEN}`,
        Accept: "application/vnd.github+json",
        "X-GitHub-Api-Version": "2022-11-28",
        "User-Agent": "colruyt-scheduler",
      },
      body: JSON.stringify({ ref: "main" }),
    });

    if (!res.ok) {
      throw new Error(`Workflow dispatch failed: ${res.status} ${await res.text()}`);
    }
  },
};
