/*
 *drbdtop - continously update stats on drbd
 *Copyright © 2017 Hayley Swimelar
 *
 *This program is free software; you can redistribute it and/or modify
 *it under the terms of the GNU General Public License as published by
 *the Free Software Foundation; either version 2 of the License, or
 *(at your option) any later version.
 *
 *This program is distributed in the hope that it will be useful,
 *but WITHOUT ANY WARRANTY; without even the implied warranty of
 *MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 *GNU General Public License for more details.
 *
 *You should have received a copy of the GNU General Public License
 *along with this program; if not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/hayswim/drbdtop/pkg/resource"
)

func main() {
	file := flag.String("file", "", "Path to a file containing output gathered from polling `drbdsetup events2 --timestamps --statistics --now`")

	flag.Parse()

	rawEvents := make(chan string)

	if *file != "" {
		f, err := os.Open(*file)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		go func() {
			for scanner.Scan() {
				rawEvents <- scanner.Text()
			}
			os.Exit(0)
		}()

	} else {
		fmt.Println("Only reading from a file is supported right now. Use the --file option next time.")
		os.Exit(1)
	}

	// Main update loop. For now just prints events.
	i := 0
	for {
		var wg sync.WaitGroup
		i++
		for {
			s := <-rawEvents
			e, err := resource.NewEvent(s)
			if err != nil {
				fmt.Printf("%v\n", err)
			}

			// Break on these event targets so that updates are applied in order.
			if e.Target == "-" {
				break
			}

			wg.Add(1)
			// Update logic goes here.
			go func() {
				defer wg.Done()
				fmt.Printf("%v\n", e)
			}()
		}
		wg.Wait()
		fmt.Printf("\nThat's %d groups!\n", i)
	}
}
