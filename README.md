Tmux Trackpad plugin
===================

# tmux plugin for magic trackpad battery(macOS only)


Features
--------
- Battery usage
- Charging status

Installation
------------
Best installed through [Tmux Plugin Manager](https://github.com/tmux-plugins/tpm) (TMP). Add following line to your `.tmux.conf` file:

```
set -g @plugin 'akayj/tmux-touchpad-battery'
```

Use `prefix + I` from inside tmux to install all plugins and source them.

Basic usage
-----------

Once plugged in, tmux `status-left` or `status-right` options can be configured with following placeholders. Each placeholder will be expanded to metric's default output.

```
set -g status-right "#{touchpad_battery}"
```

Customize options
-----------------

Available tmux options are:

- `tpb_percent_prefix` default: `TB:`
- `tpb_percent_suffix` default: `%`
- `tpb_color_charging` default: 'green'
- `tpb_color_high` default: `white`
- `tpb_color_medium` default: `yellow`
- `tpb_color_stress` default: `red`
- `tpb_stress_threshold` default: `30`
- `tpb_medium_threshold` default: `80`
- `tpb_not_show_threshold` default: `100`

