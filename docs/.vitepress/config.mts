import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Gokord Docs",
  description: "Simple & powerful Discord library",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Getting started', link: '/getting-started' },
      { text: 'Slash commands', link: '/slash-commands/' },
      { text: 'Databases', link: '/databases/' },
      { text: 'Config', link: '/config' }
    ],

    sidebar: [
      {
        text: 'Introduction',
        items: [
          { text: 'Getting started', link: '/getting-started' },
          { text: 'Custom config files', link: '/config' },
          { text: 'Innovation & Versioning', link: '/innovation' },
          { text: 'Statuses', link: 'statuses' }
        ]
      },
	  {
        text: 'Slash commands',
        items: [
          { text: 'Declaring new slash commands', link: '/slash-commands/' },
          { text: 'Using options', link: '/slash-commands/options' },
          { text: 'Using subcommands', link: '/slash-commands/sub' }
        ]
      },
	  {
        text: 'Databases',
        items: [
          { text: 'How to use databases', link: '/databases/' },
          { text: 'SQL databases', link: '/databases/sql' },
          { text: 'Redis databases', link: '/databases/redis' }
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/anhgelus/gokord/' }
    ]
  }
})
