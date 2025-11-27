---
title: "Installation"
description: "Steps to install PCT on Windows, macOS and Linux."
category: narrative
tags:
  - fundamentals
  - install
weight: 10
---

While the PCT is in early release, we provide an archive and a simple script to unpack it. When we move closer to a full release we will add a platform specific installer. Use the `install.[ps1|sh]` script, depending upon your OS:

### Bash

```bash
curl -L https://raw.githubusercontent.com/jay7x/pct/refs/heads/main/scripts/install.sh | bash
```

### PowerShell

```powershell
iex "&{ $(irm 'https://raw.githubusercontent.com/jay7x/pct/refs/heads/main/scripts/install.ps1'); Install-Pct }"
```

This will install the latest release of PCT to `~/.puppetlabs/pct`.

![install_pct](https://github.com/jay7x/pct/blob/main/docs/_resources/install_and_export_path.gif?raw=true)

> **Warning!**
>
> If you do not use the install script and are extracting the archive yourself, be sure to use the fully qualified path to `~/.puppetlabs/pct` on *nix or `$HOME/.puppetlabs/pct` on Windows when you set your `PATH` environment variable.

### Setting up Tab Completion

After installation, we'd highly recommend setting up tab completion for your shell to ensure the best possible experience.

PCT has built in tab completion support for the following shells: `bash`, `zsh`, `fish` and `powershell`.

To view the install instructions, access the `--help` menu in `pct completion` and follow the instructions for your shell:

![tab_completion](https://github.com/jay7x/pct/blob/main/docs/_resources/completion_setup.gif?raw=true)
