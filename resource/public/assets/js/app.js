/* ==========================================================================
   xyhelper-nav — public navigation page.

   Renders the catalogue returned by GET /api/v1/nav/index and provides live
   client side filtering, theming and click statistics. No build step and no
   third party runtime dependency.
   ========================================================================== */
(function () {
    'use strict';

    var API = '/api/v1';
    var THEME_KEY = 'xyhelper-nav-theme';
    var TOP_HASH = '#top';
    var EXTERNAL_ENGINES = [
        { name: '百度', url: 'https://www.baidu.com/s?wd=' },
        { name: 'Google', url: 'https://www.google.com/search?q=' },
        { name: 'Bing', url: 'https://www.bing.com/search?q=' }
    ];

    var state = {
        settings: {},
        categories: [],
        keyword: '',
        openNewTab: true,
        showHits: true
    };

    var dom = {};

    /* --- Small helpers ---------------------------------------------------- */

    function byId(id) {
        return document.getElementById(id);
    }

    function icon(name) {
        return '<svg class="icon"><use href="assets/icons.svg#' + name + '"></use></svg>';
    }

    function escapeHtml(value) {
        return String(value === null || value === undefined ? '' : value)
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;')
            .replace(/'/g, '&#39;');
    }

    function hostOf(url) {
        try {
            return new URL(url, window.location.href).host;
        } catch (error) {
            return '';
        }
    }

    /** colorFor derives a stable gradient from a site name. */
    function colorFor(text) {
        var hash = 0;
        var value = String(text || '?');
        for (var i = 0; i < value.length; i++) {
            hash = (hash * 31 + value.charCodeAt(i)) % 360;
        }
        return 'linear-gradient(135deg, hsl(' + hash + ', 64%, 54%), hsl(' + ((hash + 32) % 360) + ', 66%, 46%))';
    }

    function toast(message, type) {
        var node = document.createElement('div');
        node.className = 'toast ' + (type || '');
        node.textContent = message;
        dom.toastHost.appendChild(node);
        window.setTimeout(function () {
            node.classList.add('is-out');
            window.setTimeout(function () {
                node.remove();
            }, 260);
        }, 2600);
    }

    /** request performs a JSON API call and unwraps the GoFrame response envelope. */
    function request(path, options) {
        var config = options || {};
        config.headers = { 'Content-Type': 'application/json' };
        config.credentials = 'same-origin';
        return fetch(API + path, config).then(function (response) {
            return response.json().catch(function () {
                return null;
            }).then(function (body) {
                if (body && typeof body.code !== 'undefined' && body.code !== 0) {
                    throw new Error(body.message || '请求失败');
                }
                if (!response.ok) {
                    throw new Error('请求失败 (' + response.status + ')');
                }
                return body ? body.data : null;
            });
        });
    }

    /* --- Theme ------------------------------------------------------------ */

    function preferredTheme() {
        var stored = window.localStorage.getItem(THEME_KEY);
        if (stored === 'dark' || stored === 'light') {
            return stored;
        }
        var configured = state.settings.default_theme;
        if (configured === 'dark' || configured === 'light') {
            return configured;
        }
        return window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark';
    }

    function applyTheme(theme) {
        document.documentElement.setAttribute('data-theme', theme);
        var use = dom.themeToggle.querySelector('use');
        if (use) {
            use.setAttribute('href', 'assets/icons.svg#' + (theme === 'dark' ? 'i-sun' : 'i-moon'));
        }
    }

    /* --- Settings --------------------------------------------------------- */

    function setMeta(name, content) {
        var node = document.querySelector('meta[name="' + name + '"]');
        if (node && content) {
            node.setAttribute('content', content);
        }
    }

    function applySettings() {
        var settings = state.settings;
        var siteName = settings.site_name || '导航';

        document.title = siteName;
        setMeta('description', settings.site_description);
        setMeta('keywords', settings.site_keywords);

        dom.heroTitle.textContent = siteName;
        dom.brandName.textContent = siteName;
        dom.brandLogo.textContent = '★';

        if (settings.site_subtitle) {
            dom.heroSubtitle.textContent = settings.site_subtitle;
        } else {
            dom.heroSubtitle.remove();
        }

        if (settings.theme_color) {
            document.documentElement.style.setProperty('--accent', settings.theme_color);
        }

        if (settings.announcement) {
            dom.announcement.textContent = settings.announcement;
            dom.announcement.classList.add('is-visible');
        }

        if (settings.background_image && settings.background_image.indexOf('"') === -1) {
            dom.ambient.style.backgroundImage =
                'linear-gradient(rgba(7, 11, 24, 0.74), rgba(7, 11, 24, 0.86)), url("' +
                settings.background_image + '")';
            dom.ambient.style.backgroundSize = 'cover';
            dom.ambient.style.backgroundPosition = 'center';
        }

        var footer = [];
        if (settings.site_footer) {
            footer.push(escapeHtml(settings.site_footer));
        }
        if (settings.icp) {
            footer.push('<a href="https://beian.miit.gov.cn" target="_blank" rel="noopener">' +
                escapeHtml(settings.icp) + '</a>');
        }
        if (footer.length === 0) {
            footer.push('Powered by xyhelper-nav');
        }
        dom.footer.innerHTML = footer.join('<span class="sep">·</span>');
    }

    /* --- Rendering -------------------------------------------------------- */

    function categoryEmoji(category) {
        var value = (category.icon || '').trim();
        if (!value) {
            return '📁';
        }
        return value.indexOf('http') === 0 || value.indexOf('/') === 0 ? '📁' : value;
    }

    /**
     * siteIconHtml points the image at the local icon cache. The server
     * downloads a remote icon once and serves it from this origin afterwards,
     * so the browser never issues a request to a third party host. When no icon
     * is available the image removes itself and the generated letter avatar
     * underneath stays visible.
     */
    function siteIconHtml(site) {
        if (!site.iconUrl) {
            return '';
        }
        return '<img src="' + escapeHtml(site.iconUrl) + '" alt="" loading="lazy" decoding="async"' +
            ' onerror="this.remove()">';
    }

    function cardHtml(site, index) {
        var host = hostOf(site.url);
        var letter = (site.name || host || '?').trim().charAt(0);
        var target = state.openNewTab ? ' target="_blank" rel="noopener noreferrer"' : '';
        // Showing a zero counter on every card is noise, so only display it once
        // a site has actually been visited.
        var hits = state.showHits && site.hits > 0
            ? '<span class="hits" data-hits>' + icon('i-chart') + '<span>' + site.hits + '</span></span>'
            : '';
        var description = site.description
            ? '<div class="card-desc">' + escapeHtml(site.description) + '</div>'
            : '';
        var hostLine = host ? '<div class="card-host">' + escapeHtml(host) + '</div>' : '';
        var delay = Math.min(index * 22, 320);

        return '<a class="card" href="' + escapeHtml(site.url) + '" data-id="' + site.id + '"' +
            ' data-name="' + escapeHtml(site.name) + '"' + target +
            ' style="animation-delay:' + delay + 'ms">' +
            '<span class="card-icon" style="background-image:' + colorFor(site.name) + '">' +
            '<span class="letter">' + escapeHtml(letter) + '</span>' + siteIconHtml(site) +
            '</span>' +
            '<span class="card-body">' +
            '<span class="card-title"><span class="name">' + escapeHtml(site.name) + '</span>' + hits + '</span>' +
            description + hostLine +
            '</span></a>';
    }

    function emptyState(title, hint) {
        return '<div class="empty">' + icon('i-inbox') +
            '<div>' + escapeHtml(title) + '</div>' +
            (hint ? '<div style="margin-top:6px;font-size:13px">' + escapeHtml(hint) + '</div>' : '') +
            '</div>';
    }

    function renderSidebar() {
        var total = 0;
        state.categories.forEach(function (category) {
            total += category.sites.length;
        });
        var html = '<a class="nav-item active" href="' + TOP_HASH + '" data-target="all">' +
            '<span class="emoji">★</span><span>全部</span>' +
            '<span class="count">' + total + '</span></a>';
        state.categories.forEach(function (category) {
            html += '<a class="nav-item" href="#cat-' + category.id + '" data-target="cat-' + category.id + '">' +
                '<span class="emoji">' + escapeHtml(categoryEmoji(category)) + '</span>' +
                '<span>' + escapeHtml(category.name) + '</span>' +
                '<span class="count">' + category.sites.length + '</span></a>';
        });
        dom.navList.innerHTML = html;
    }

    function renderContent() {
        if (state.categories.length === 0) {
            dom.content.innerHTML = emptyState('还没有添加任何站点', '登录管理后台即可添加分类与站点');
            return;
        }
        var html = '';
        state.categories.forEach(function (category) {
            if (category.sites.length === 0) {
                return;
            }
            var description = category.description
                ? '<span class="cat-desc" style="color:var(--text-faint);font-size:12.5px">' +
                  escapeHtml(category.description) + '</span>'
                : '';
            html += '<section class="cat" id="cat-' + category.id + '">' +
                '<div class="cat-head">' +
                '<span class="emoji">' + escapeHtml(categoryEmoji(category)) + '</span>' +
                '<h2>' + escapeHtml(category.name) + '</h2>' +
                '<span class="count">' + category.sites.length + '</span>' +
                description +
                '<span class="line"></span></div>' +
                '<div class="grid">' +
                category.sites.map(function (site, index) {
                    return cardHtml(site, index);
                }).join('') +
                '</div></section>';
        });
        if (html === '') {
            html = emptyState('还没有添加任何站点', '登录管理后台即可添加分类与站点');
        }
        dom.content.innerHTML = html + '<div class="search-empty" id="searchEmpty"></div>';
    }

    /** markText renders text into an element, wrapping keyword matches in <mark>. */
    function markText(element, text, keyword) {
        if (!keyword) {
            element.textContent = text;
            return;
        }
        var lower = text.toLowerCase();
        var needle = keyword.toLowerCase();
        var html = '';
        var cursor = 0;
        var found = lower.indexOf(needle, cursor);
        if (found === -1) {
            element.textContent = text;
            return;
        }
        while (found !== -1) {
            html += escapeHtml(text.slice(cursor, found)) +
                '<mark>' + escapeHtml(text.slice(found, found + needle.length)) + '</mark>';
            cursor = found + needle.length;
            found = lower.indexOf(needle, cursor);
        }
        element.innerHTML = html + escapeHtml(text.slice(cursor));
    }

    function renderEngines() {
        var keyword = state.keyword.trim();
        if (!keyword) {
            dom.engines.innerHTML = '';
            return;
        }
        dom.engines.innerHTML = EXTERNAL_ENGINES.map(function (engine) {
            return '<a class="engine" target="_blank" rel="noopener" href="' + engine.url +
                encodeURIComponent(keyword) + '">' + engine.name + '</a>';
        }).join('');
    }

    /** applyFilter hides non matching cards and highlights the matches. */
    function applyFilter(keyword) {
        state.keyword = keyword;
        var needle = keyword.trim().toLowerCase();
        var total = 0;

        state.categories.forEach(function (category) {
            var section = byId('cat-' + category.id);
            if (!section) {
                return;
            }
            var visible = 0;
            category.sites.forEach(function (site) {
                var card = section.querySelector('.card[data-id="' + site.id + '"]');
                if (!card) {
                    return;
                }
                var haystack = (site.name + ' ' + site.description + ' ' + site.url + ' ' + category.name).toLowerCase();
                var matched = needle === '' || haystack.indexOf(needle) !== -1;
                card.classList.toggle('is-hidden', !matched);
                markText(card.querySelector('.name'), site.name, matched ? needle : '');
                if (matched) {
                    visible++;
                }
            });
            section.classList.toggle('is-hidden', visible === 0);
            var badge = section.querySelector('.cat-head .count');
            if (badge) {
                badge.textContent = visible;
            }
            var navCount = dom.navList.querySelector('.nav-item[data-target="cat-' + category.id + '"] .count');
            if (navCount) {
                navCount.textContent = visible;
            }
            total += visible;
        });

        var allCount = dom.navList.querySelector('.nav-item[data-target="all"] .count');
        if (allCount) {
            allCount.textContent = total;
        }
        dom.searchStat.textContent = needle === ''
            ? ('共 ' + total + ' 个站点')
            : ('匹配 ' + total + ' 个站点');
        renderEngines();

        var empty = byId('searchEmpty');
        if (empty) {
            if (needle !== '' && total === 0) {
                empty.innerHTML = '<div>' + icon('i-search') + '</div>' +
                    '<p>没有找到与 “<strong>' + escapeHtml(keyword.trim()) + '</strong>” 匹配的站点</p>' +
                    '<p style="margin-top:8px;font-size:13px">试试在这些搜索引擎中查找</p>' +
                    '<div class="engines" style="margin-top:12px;justify-content:center">' +
                    EXTERNAL_ENGINES.map(function (engine) {
                        return '<a class="engine" target="_blank" rel="noopener" href="' + engine.url +
                            encodeURIComponent(keyword.trim()) + '">' + engine.name + '</a>';
                    }).join('') + '</div>';
                empty.classList.add('is-visible');
            } else {
                empty.classList.remove('is-visible');
                empty.innerHTML = '';
            }
        }

        var content = dom.content.querySelector('.empty');
        if (content) {
            content.classList.toggle('is-hidden', needle !== '' && total > 0);
        }
    }

    /* --- Interactions ----------------------------------------------------- */

    function recordHit(id) {
        if (!id) {
            return;
        }
        request('/nav/hit', {
            method: 'POST',
            body: JSON.stringify({ id: Number(id) }),
            keepalive: true
        }).then(function (data) {
            if (!data || !state.showHits) {
                return;
            }
            var badge = dom.content.querySelector('.card[data-id="' + id + '"] [data-hits] span');
            if (badge) {
                badge.textContent = data.hits;
            }
        }).catch(function () {
            /* Counting a click must never disturb navigation. */
        });
    }

    function observeSections() {
        var sections = dom.content.querySelectorAll('.cat');
        if (!('IntersectionObserver' in window) || sections.length === 0) {
            return;
        }
        var observer = new IntersectionObserver(function (entries) {
            if (state.keyword.trim() !== '') {
                return;
            }
            entries.forEach(function (entry) {
                if (!entry.isIntersecting) {
                    return;
                }
                var id = entry.target.getAttribute('id');
                dom.navList.querySelectorAll('.nav-item').forEach(function (item) {
                    item.classList.toggle('active', item.getAttribute('data-target') === id);
                });
            });
        }, { rootMargin: '-96px 0px -70% 0px', threshold: 0 });
        sections.forEach(function (section) {
            observer.observe(section);
        });
    }

    function isTyping(target) {
        if (!target) {
            return false;
        }
        var tag = target.tagName;
        return tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT' || target.isContentEditable;
    }

    function bindEvents() {
        dom.themeToggle.addEventListener('click', function () {
            var next = document.documentElement.getAttribute('data-theme') === 'dark' ? 'light' : 'dark';
            window.localStorage.setItem(THEME_KEY, next);
            applyTheme(next);
        });

        dom.searchInput.addEventListener('input', function () {
            applyFilter(dom.searchInput.value);
        });

        dom.searchInput.addEventListener('keydown', function (event) {
            if (event.key === 'Enter') {
                var first = dom.content.querySelector('.card:not(.is-hidden)');
                if (first) {
                    window.open(first.href, state.openNewTab ? '_blank' : '_self');
                }
            } else if (event.key === 'Escape') {
                dom.searchInput.value = '';
                applyFilter('');
                dom.searchInput.blur();
            }
        });

        document.addEventListener('keydown', function (event) {
            if (isTyping(event.target)) {
                return;
            }
            if (event.key === '/' || ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === 'k')) {
                event.preventDefault();
                dom.searchInput.focus();
                dom.searchInput.select();
            }
        });

        document.addEventListener('click', function (event) {
            var card = event.target.closest ? event.target.closest('.card') : null;
            if (card) {
                recordHit(card.getAttribute('data-id'));
            }
        });

        window.addEventListener('scroll', function () {
            dom.toTop.classList.toggle('is-visible', window.scrollY > 420);
        }, { passive: true });

        dom.toTop.addEventListener('click', function () {
            window.scrollTo({ top: 0, behavior: 'smooth' });
        });
    }

    function init() {
        dom = {
            ambient: byId('ambient'),
            themeToggle: byId('themeToggle'),
            heroTitle: byId('heroTitle'),
            heroSubtitle: byId('heroSubtitle'),
            brandName: byId('brandName'),
            brandLogo: byId('brandLogo'),
            announcement: byId('announcement'),
            searchInput: byId('searchInput'),
            searchStat: byId('searchStat'),
            engines: byId('engines'),
            navList: byId('navList'),
            content: byId('content'),
            footer: byId('footer'),
            toTop: byId('toTop'),
            toastHost: byId('toastHost')
        };

        applyTheme(preferredTheme());
        bindEvents();

        request('/nav/index').then(function (data) {
            state.settings = (data && data.settings) || {};
            state.categories = (data && data.categories) || [];
            state.openNewTab = state.settings.open_new_tab !== '0';
            state.showHits = state.settings.show_hits !== '0';

            applySettings();
            applyTheme(preferredTheme());
            renderSidebar();
            renderContent();
            applyFilter('');
            observeSections();
        }).catch(function (error) {
            toast(error.message, 'error');
            dom.content.innerHTML = emptyState('加载失败', error.message);
        });
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
