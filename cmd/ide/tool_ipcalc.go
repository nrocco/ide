package main

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/fatih/color"
	"github.com/nrocco/ide/pkg/ide/tools"
	"github.com/spf13/cobra"
)

// ipcalcMaxSubnets is the maximum number of subnets listed with --split
const ipcalcMaxSubnets = 64

var ipcalcOpts struct {
	Split    int
	Contains string
}

var ipcalcCmd = &cobra.Command{
	Use:   "ipcalc cidr",
	Short: "Show details for an IPv4 or IPv6 CIDR range",
	Long:  "Show details for an IPv4 or IPv6 CIDR range, e.g. 192.168.1.10/24 or 2001:db8::/112 (host bits allowed)",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		network, err := tools.NewIPNetwork(args[0])
		if err != nil {
			return err
		}

		address := network.Address()
		netmask := network.Netmask()
		out := cmd.OutOrStdout()

		blue := color.New(color.FgBlue).SprintFunc()
		grey := color.New(color.FgHiBlack).SprintfFunc()

		// IPv6 addresses also show their exploded notation
		exploded := func(ip tools.IPAddress) string {
			if network.Version() == 4 {
				return ""
			}
			return " " + grey("(%s)", ip.Exploded())
		}

		fmt.Fprintf(out, "Version                    : %s\n", blue(fmt.Sprintf("IPv%d", network.Version())))
		if network.Version() == 4 {
			fmt.Fprintf(out, "Address                    : %s %s\n", blue(address), grey("(%s)", address.Binary()))
		} else {
			fmt.Fprintf(out, "Address                    : %s%s\n", blue(address), exploded(address))
		}
		fmt.Fprintf(out, "Is Network Address         : %s\n", blue(address == network.NetworkAddress()))
		fmt.Fprintf(out, "Network                    : %s\n", blue(network))
		fmt.Fprintf(out, "Prefix Length              : %s\n", blue(network.PrefixLength()))
		if network.Version() == 4 {
			fmt.Fprintf(out, "Netmask                    : %s %s\n", blue(netmask), grey("(%s) (%s)", netmask.Hex(), netmask.Binary()))
		} else {
			fmt.Fprintf(out, "Netmask                    : %s\n", blue(netmask))
		}
		fmt.Fprintf(out, "Hostmask Wildcard          : %s\n", blue(network.Hostmask()))
		fmt.Fprintf(out, "Network Address            : %s%s\n", blue(network.NetworkAddress()), exploded(network.NetworkAddress()))
		if network.Version() == 4 {
			if broadcast, ok := network.Broadcast(); ok {
				fmt.Fprintf(out, "Broadcast Address          : %s\n", blue(broadcast))
			} else {
				fmt.Fprintf(out, "Broadcast Address          : %s %s\n", blue("none"), grey("(/31 or /32)"))
			}
		} else {
			fmt.Fprintf(out, "Last Address               : %s%s\n", blue(network.LastAddress()), exploded(network.LastAddress()))
		}
		fmt.Fprintf(out, "First Usable               : %s%s\n", blue(network.FirstUsable()), exploded(network.FirstUsable()))
		fmt.Fprintf(out, "Last Usable                : %s%s\n", blue(network.LastUsable()), exploded(network.LastUsable()))
		fmt.Fprintf(out, "Usable / Total Addresses   : %s %s\n", blue(formatThousands(network.NumUsable())+" / "+formatThousands(network.NumAddresses())), grey("(%s)", network.UsableNote()))
		fmt.Fprintf(out, "Flags                      : %s\n", blue(formatList(network.Flags())))

		if network.Version() == 6 {
			if ipv4, ok := address.IPv4Mapped(); ok {
				fmt.Fprintf(out, "IPv4 Mapped                : %s\n", blue(ipv4))
			}
			if ipv4, ok := address.SixToFour(); ok {
				fmt.Fprintf(out, "6to4 IPv4                  : %s\n", blue(ipv4))
			}
			if server, client, ok := address.Teredo(); ok {
				fmt.Fprintf(out, "Teredo Server Client       : %s\n", blue(fmt.Sprintf("%s, %s", server, client)))
			}
		}

		// Neighbouring blocks of the same size and the parent block
		fmt.Fprintf(out, "Previous Network           : %s\n", blue(formatOptional(network.Previous())))
		fmt.Fprintf(out, "Next Network               : %s\n", blue(formatOptional(network.Next())))
		fmt.Fprintf(out, "Supernet                   : %s\n", blue(formatOptional(network.Supernet())))

		if ipcalcOpts.Contains != "" {
			other, err := tools.NewIPNetwork(ipcalcOpts.Contains)
			if err != nil {
				return fmt.Errorf("--contains: %w", err)
			}
			fmt.Fprintf(out, "Contains                   : %s\n", blue(fmt.Sprintf("%s contained=%t overlaps=%t", ipcalcOpts.Contains, network.Contains(other), network.Overlaps(other))))
		}

		if cmd.Flags().Changed("split") {
			subnets, err := network.Subnets(ipcalcOpts.Split)
			if err != nil {
				return fmt.Errorf("--split: %w", err)
			}

			var lines []string
			var size *big.Int
			for subnet := range subnets {
				if len(lines) == ipcalcMaxSubnets {
					lines = append(lines, fmt.Sprintf("... (first %d shown)", ipcalcMaxSubnets))
					break
				}
				size = subnet.NumAddresses()
				lines = append(lines, subnet.String())
			}
			count := new(big.Int).Div(network.NumAddresses(), size)

			fmt.Fprintf(out, "Split                      : %s\n", blue(fmt.Sprintf("/%d -> %s subnets, %s addresses each", ipcalcOpts.Split, count, size)))
			for _, line := range lines {
				fmt.Fprintf(out, "                             %s\n", blue(line))
			}
		}

		return nil
	},
}

// formatOptional returns the network, or - when it is not available
func formatOptional(network tools.IPNetwork, ok bool) string {
	if !ok {
		return "-"
	}
	return network.String()
}

func formatList(values []string) string {
	if len(values) == 0 {
		return "-"
	}
	return strings.Join(values, ", ")
}

func formatThousands(n *big.Int) string {
	s := n.String()
	sign := ""
	if strings.HasPrefix(s, "-") {
		sign, s = "-", s[1:]
	}
	var b strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(c)
	}
	return sign + b.String()
}

func init() {
	ipcalcCmd.Flags().IntVar(&ipcalcOpts.Split, "split", 0, "split range into subnets of this prefix length")
	ipcalcCmd.Flags().StringVar(&ipcalcOpts.Contains, "contains", "", "check if address/range is inside and overlaps")

	toolCmd.AddCommand(ipcalcCmd)
}
