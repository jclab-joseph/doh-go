/*
 * Copyright 2019 Li Kexian
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * DNS over HTTPS (DoH) Golang implementation
 * https://www.likexian.com/
 */

package google

import (
	"context"
	"github.com/jclab-joseph/doh-go/dns"
	"github.com/likexian/gokit/assert"
	"net/http"
	"testing"
	"time"
)

func TestVersion(t *testing.T) {
	assert.Contains(t, Version(), ".")
	assert.Contains(t, Author(), "likexian")
	assert.Contains(t, License(), "Apache License")
}

func TestString(t *testing.T) {
	c := New()
	assert.Equal(t, c.String(), "google")
}

func TestQuery(t *testing.T) {
	c := New()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rsp, err := c.Query(ctx, http.DefaultClient, "likexian.com", dns.TypeA)
	assert.Nil(t, err)
	assert.Gt(t, len(rsp.Answer), 0)

	rsp, err = c.Query(ctx, http.DefaultClient, "한국인터넷정보센터.한국", dns.TypeA)
	assert.Nil(t, err)
	assert.Gt(t, len(rsp.Answer), 0)
}

func TestECSQuery(t *testing.T) {
	c := New()
	err := c.SetProvides(DefaultProvides)
	assert.Nil(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = c.ECSQuery(ctx, http.DefaultClient, "xx", dns.TypeA, "1.1.1.1")
	assert.NotNil(t, err)

	_, err = c.ECSQuery(ctx, http.DefaultClient, "likexian.com", dns.TypeA, "xx")
	assert.NotNil(t, err)

	rsp, err := c.ECSQuery(ctx, http.DefaultClient, "likexian.com", dns.TypeA, "1.1.1.1")
	assert.Nil(t, err)
	assert.Gt(t, len(rsp.Answer), 0)

	rsp, err = c.ECSQuery(ctx, http.DefaultClient, "likexian.com", dns.TypeA, "1.1.1.1/24")
	assert.Nil(t, err)
	assert.Gt(t, len(rsp.Answer), 0)

	Upstream[DefaultProvides] = "test"
	_, err = c.ECSQuery(ctx, http.DefaultClient, "likexian.com", dns.TypeA, "")
	assert.NotNil(t, err)

	Upstream[DefaultProvides] = "https://dns.google.com/dns"
	_, err = c.ECSQuery(ctx, http.DefaultClient, "likexian.com", dns.TypeA, "")
	assert.NotNil(t, err)
}
