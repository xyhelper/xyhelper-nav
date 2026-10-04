/* ==========================================================================
   xyhelper-nav — administration panel.

   Talks to the protected /api/v1/admin/* endpoints using the session cookie
   issued by the login endpoint. Plain JavaScript, no build step.
   ========================================================================== */
(function () {
    'use strict';

    var API = '/api/v1';
    var THEME_KEY = 'xyhelper-nav-theme';
    var VIEW_TITLES = {
        sites: '站点管理',
        categories: '分类管理',
        settings: '站点设置',
        data: '数据备份',
        account: '账号安全'
    };

    /**
     * SETTING_SECTIONS drives the dynamically generated settings form.
     * Every entry maps to one key of the nav_setting table.
     */
    var SETTING_SECTIONS = [
        {
            target: 'settingsBasic',
            fields: [
                { key: 'site_name', label: '站点名称', type: 'text', placeholder: 'XYHelper 导航' },
                { key: 'site_subtitle', label: '副标题', type: 'text', placeholder: '轻量 · 美观 · 自托管导航页' },
                { key: 'site_description', label: '站点描述', type: 'textarea', hint: '用于搜索引擎收录' },
                { key: 'site_keywords', label: '站点关键词', type: 'text', hint: '多个关键词以英文逗号分隔' },
                { key: 'announcement', label: '首页公告', type: 'textarea', hint: '留空则不显示公告栏' }
            ]
        },
        {
            target: 'settingsAppearance',
            fields: [
                { key: 'theme_color', label: '主题色', type: 'color' },
                {
                    key: 'default_theme',
                    label: '默认主题',
                    type: 'select',
                    options: [['dark', '深色'], ['light', '浅色']]
                },
                { key: 'background_image', label: '背景图地址', type: 'text', placeholder: 'https://…' }
            ]
        },
        {
            target: 'settingsBehaviour',
            fields: [
                { key: 'open_new_tab', label: '新标签页打开站点', type: 'switch' },
                { key: 'show_hits', label: '显示点击次数', type: 'switch' }
            ]
        },
        {
            target: 'settingsFooter',
            fields: [
                { key: 'site_footer', label: '页脚文案', type: 'text' },
                { key: 'icp', label: 'ICP 备案号', type: 'text', placeholder: '京ICP备00000000号' }
            ]
        }
    ];

    var state = {
        categories: [],
        settings: {},
        profile: null,
        view: 'sites',
        filter: { categoryId: 0, keyword: '' }
    };

    var modal = { confirm: null, cancel: null };
    var dom = {};

    /* --- Helpers ---------------------------------------------------------- */

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

    /** colorFor derives a stable avatar gradient from a name. */
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

    /** request performs an API call and unwraps the GoFrame response envelope. */
    function request(path, options) {
        var config = options || {};
        config.headers = { 'Content-Type': 'application/json' };
        config.credentials = 'same-origin';
        return fetch(API + path, config).then(function (response) {
            return response.json().catch(function () {
                return null;
            }).then(function (body) {
                if (response.status === 401) {
                    var unauthorized = new Error('登录状态已失效，请重新登录');
                    unauthorized.unauthorized = true;
                    throw unauthorized;
                }
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

    function handleError(error) {
        if (error && error.unauthorized) {
            showLogin();
            return;
        }
        toast(error.message || '操作失败', 'error');
    }

    /* --- Theme ------------------------------------------------------------ */

    function applyTheme(theme) {
        document.documentElement.setAttribute('data-theme', theme);
        var use = dom.themeToggle ? dom.themeToggle.querySelector('use') : null;
        if (use) {
            use.setAttribute('href', 'assets/icons.svg#' + (theme === 'dark' ? 'i-sun' : 'i-moon'));
        }
    }

    function initialTheme() {
        var stored = window.localStorage.getItem(THEME_KEY);
        if (stored === 'dark' || stored === 'light') {
            return stored;
        }
        return window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark';
    }

    /* --- Modal ------------------------------------------------------------ */

    function openModal(options) {
        modal.confirm = options.onConfirm || null;
        modal.cancel = options.onCancel || null;
        dom.modalTitle.textContent = options.title || '';
        dom.modalBody.innerHTML = options.bodyHtml || '';
        dom.modalConfirm.textContent = options.confirmText || '确定';
        dom.modalConfirm.className = 'btn ' + (options.danger ? 'btn-danger' : 'btn-primary');
        dom.modalMask.classList.add('is-open');
        var firstInput = dom.modalBody.querySelector('input, textarea, select');
        if (firstInput) {
            firstInput.focus();
        }
    }

    function closeModal() {
        dom.modalMask.classList.remove('is-open');
        dom.modalBody.innerHTML = '';
        modal.confirm = null;
        modal.cancel = null;
    }

    function cancelModal() {
        var onCancel = modal.cancel;
        closeModal();
        if (onCancel) {
            onCancel();
        }
    }

    /** confirmDialog shows a modal confirmation and resolves to a boolean. */
    function confirmDialog(title, message, options) {
        var config = options || {};
        return new Promise(function (resolve) {
            openModal({
                title: title,
                bodyHtml: '<p style="color:var(--text-dim);font-size:14px">' + escapeHtml(message) + '</p>',
                confirmText: config.confirmText || '确定',
                danger: config.danger,
                onConfirm: function () {
                    resolve(true);
                    return true;
                },
                onCancel: function () {
                    resolve(false);
                }
            });
        });
    }

    /* --- Views ------------------------------------------------------------ */

    function switchView(view) {
        state.view = view;
        dom.viewTitle.textContent = VIEW_TITLES[view] || '';
        Array.prototype.forEach.call(dom.adminNav.querySelectorAll('button'), function (button) {
            button.classList.toggle('active', button.getAttribute('data-view') === view);
        });
        ['sites', 'categories', 'settings', 'data', 'account'].forEach(function (name) {
            byId('panel-' + name).classList.toggle('active', name === view);
        });
        dom.viewSubtitle.textContent = {
            sites: '维护导航页上的站点卡片',
            categories: '分类决定站点在导航页上的分组',
            settings: '这些配置会立即应用到导航页',
            data: '导出备份或从备份恢复数据',
            account: '管理登录账号与密码'
        }[view] || '';
    }

    function showLogin() {
        dom.adminView.hidden = true;
        dom.loginView.hidden = false;
        dom.loginPassword.value = '';
    }

    function showAdmin() {
        dom.loginView.hidden = true;
        dom.adminView.hidden = false;
    }

    /* --- Loading ---------------------------------------------------------- */

    function loadCategories() {
        return request('/admin/categories').then(function (data) {
            state.categories = (data && data.list) || [];
        });
    }

    function loadSettings() {
        return request('/admin/settings').then(function (data) {
            state.settings = (data && data.settings) || {};
            applyAccent();
            // Follow the configured default theme unless the operator has
            // explicitly picked one in this browser.
            var stored = window.localStorage.getItem(THEME_KEY);
            if (stored !== 'dark' && stored !== 'light') {
                var configured = state.settings.default_theme;
                if (configured === 'dark' || configured === 'light') {
                    applyTheme(configured);
                }
            }
            renderSettingsForm();
        });
    }

    function loadProfile() {
        return request('/admin/profile').then(function (data) {
            state.profile = data;
            renderProfile();
        });
    }

    function applyAccent() {
        if (state.settings.theme_color) {
            document.documentElement.style.setProperty('--accent', state.settings.theme_color);
        }
    }

    /* --- Rendering -------------------------------------------------------- */

    function renderProfile() {
        if (!state.profile) {
            return;
        }
        dom.profileNickname.value = state.profile.nickname || '';
        var meta = ['账号：' + state.profile.username];
        if (state.profile.lastLoginAt && state.profile.lastLoginAt.indexOf('0001-') !== 0) {
            meta.push('上次登录：' + new Date(state.profile.lastLoginAt).toLocaleString());
        }
        if (state.profile.lastLoginIp) {
            meta.push('来源 IP：' + state.profile.lastLoginIp);
        }
        dom.accountMeta.textContent = meta.join(' · ');
    }

    function renderBrand() {
        var name = state.settings.site_name || 'xyhelper-nav';
        dom.sideName.textContent = name;
        dom.sideLogo.textContent = '★';
        document.title = '管理后台 · ' + name;
    }

    function allSites() {
        var rows = [];
        state.categories.forEach(function (category) {
            (category.sites || []).forEach(function (site) {
                rows.push({ site: site, category: category });
            });
        });
        return rows;
    }

    function renderStats() {
        var rows = allSites();
        var hits = rows.reduce(function (sum, row) {
            return sum + (row.site.hits || 0);
        }, 0);
        var hidden = rows.filter(function (row) {
            return row.site.status !== 1;
        }).length;
        dom.statGrid.innerHTML = [
            { value: rows.length, label: '站点总数' },
            { value: state.categories.length, label: '分类数量' },
            { value: hits, label: '累计点击' },
            { value: hidden, label: '隐藏站点' }
        ].map(function (item) {
            return '<div class="stat"><div class="value">' + item.value + '</div>' +
                '<div class="label">' + item.label + '</div></div>';
        }).join('');
    }

    function renderCategoryFilter() {
        var current = state.filter.categoryId;
        var options = ['<option value="0">全部分类</option>'];
        state.categories.forEach(function (category) {
            options.push('<option value="' + category.id + '">' + escapeHtml(category.name) +
                '（' + (category.sites || []).length + '）</option>');
        });
        dom.siteCategoryFilter.innerHTML = options.join('');
        dom.siteCategoryFilter.value = String(current);
    }

    function siteRow(row) {
        var site = row.site;
        var letter = (site.name || '?').trim().charAt(0);
        var image = site.iconUrl
            ? '<img src="' + escapeHtml(site.iconUrl) + '" alt="" style="width:20px;height:20px;object-fit:contain"' +
              ' onerror="this.remove()">'
            : '';
        var status = site.status === 1
            ? '<span class="badge on">显示</span>'
            : '<span class="badge off">隐藏</span>';
        return '<tr>' +
            '<td><div class="cell-site">' +
            '<span class="cell-icon" style="background-image:' + colorFor(site.name) + '">' +
            escapeHtml(letter) + image + '</span>' +
            '<span><span class="cell-name">' + escapeHtml(site.name) + '</span>' +
            (site.description ? '<div class="cell-sub">' + escapeHtml(site.description) + '</div>' : '') +
            '</span></div></td>' +
            '<td class="mono">' + escapeHtml(site.url) + '</td>' +
            '<td class="nowrap">' + escapeHtml(row.category.name) + '</td>' +
            '<td>' + site.sort + '</td>' +
            '<td>' + status + '</td>' +
            '<td>' + (site.hits || 0) + '</td>' +
            '<td><div class="row-actions">' +
            '<a class="btn btn-sm btn-icon btn-ghost" href="' + escapeHtml(site.url) + '" target="_blank" rel="noopener" title="访问">' +
            icon('i-external') + '</a>' +
            '<button class="btn btn-sm btn-icon btn-ghost" type="button" data-action="edit-site" data-id="' + site.id + '" title="编辑">' +
            icon('i-edit') + '</button>' +
            '<button class="btn btn-sm btn-icon btn-ghost" type="button" data-action="delete-site" data-id="' + site.id + '" title="删除">' +
            icon('i-trash') + '</button>' +
            '</div></td></tr>';
    }

    function renderSites() {
        var keyword = state.filter.keyword.trim().toLowerCase();
        var rows = allSites().filter(function (row) {
            if (state.filter.categoryId && row.category.id !== state.filter.categoryId) {
                return false;
            }
            if (!keyword) {
                return true;
            }
            return (row.site.name + ' ' + row.site.url + ' ' + (row.site.description || ''))
                .toLowerCase().indexOf(keyword) !== -1;
        });
        dom.siteTableBody.innerHTML = rows.length
            ? rows.map(siteRow).join('')
            : '<tr><td colspan="7"><div class="empty">' + icon('i-inbox') +
              '<div>没有匹配的站点</div></div></td></tr>';
        renderStats();
    }

    function renderCategories() {
        dom.categoryTableBody.innerHTML = state.categories.length
            ? state.categories.map(function (category) {
                var status = category.status === 1
                    ? '<span class="badge on">显示</span>'
                    : '<span class="badge off">隐藏</span>';
                var visibility = category.visibility === 1
                    ? '<span class="badge warn" title="需要登录管理后台才能看到">' + icon('i-lock') + '登录可见</span>'
                    : '<span class="badge off">公开</span>';
                return '<tr>' +
                    '<td><div class="cell-site"><span class="cell-icon" style="background-image:' + colorFor(category.name) + '">' +
                    escapeHtml((category.icon || category.name || '?').trim().charAt(0)) + '</span>' +
                    '<span class="cell-name">' + escapeHtml(category.name) + '</span></div></td>' +
                    '<td class="cell-sub">' + escapeHtml(category.description || '—') + '</td>' +
                    '<td>' + (category.sites || []).length + '</td>' +
                    '<td>' + category.sort + '</td>' +
                    '<td>' + status + '</td>' +
                    '<td>' + visibility + '</td>' +
                    '<td><div class="row-actions">' +
                    '<button class="btn btn-sm btn-icon btn-ghost" type="button" data-action="edit-category" data-id="' + category.id + '" title="编辑">' +
                    icon('i-edit') + '</button>' +
                    '<button class="btn btn-sm btn-icon btn-ghost" type="button" data-action="delete-category" data-id="' + category.id + '" title="删除">' +
                    icon('i-trash') + '</button>' +
                    '</div></td></tr>';
            }).join('')
            : '<tr><td colspan="7"><div class="empty">' + icon('i-folder') +
              '<div>还没有分类，先创建一个吧</div></div></td></tr>';
        renderCategoryFilter();
    }

    function renderSettingsForm() {
        SETTING_SECTIONS.forEach(function (section) {
            var host = byId(section.target);
            host.innerHTML = section.fields.map(function (field) {
                var value = state.settings[field.key] === undefined ? '' : state.settings[field.key];
                var control;
                if (field.type === 'textarea') {
                    control = '<textarea class="textarea" id="setting_' + field.key + '">' +
                        escapeHtml(value) + '</textarea>';
                } else if (field.type === 'select') {
                    control = '<select class="select" id="setting_' + field.key + '">' +
                        field.options.map(function (option) {
                            return '<option value="' + option[0] + '"' +
                                (value === option[0] ? ' selected' : '') + '>' + option[1] + '</option>';
                        }).join('') + '</select>';
                } else if (field.type === 'switch') {
                    control = '<label class="switch"><input type="checkbox" id="setting_' + field.key + '"' +
                        (value === '1' ? ' checked' : '') + '><span class="track"></span>' +
                        '<span class="label">' + escapeHtml(field.label) + '</span></label>';
                } else if (field.type === 'color') {
                    control = '<div class="color-row"><input type="color" id="setting_' + field.key + '" value="' +
                        escapeHtml(value || '#6366f1') + '">' +
                        '<input class="input" id="setting_' + field.key + '_text" value="' + escapeHtml(value) + '"></div>';
                } else {
                    control = '<input class="input" id="setting_' + field.key + '" value="' + escapeHtml(value) + '"' +
                        (field.placeholder ? ' placeholder="' + escapeHtml(field.placeholder) + '"' : '') + '>';
                }
                return '<div class="field' + (field.type === 'textarea' ? ' full' : '') + '">' +
                    (field.type === 'switch' ? '' : '<label for="setting_' + field.key + '">' + escapeHtml(field.label) + '</label>') +
                    control +
                    (field.hint ? '<span class="hint">' + escapeHtml(field.hint) + '</span>' : '') +
                    '</div>';
            }).join('');
        });
        // Keep the colour picker and its text twin in sync.
        document.querySelectorAll('input[type="color"]').forEach(function (picker) {
            var text = byId(picker.id + '_text');
            if (!text) {
                return;
            }
            picker.addEventListener('input', function () {
                text.value = picker.value;
            });
            text.addEventListener('input', function () {
                if (/^#[0-9a-fA-F]{6}$/.test(text.value.trim())) {
                    picker.value = text.value.trim();
                }
            });
        });
        renderBrand();
    }

    function collectSettings() {
        var payload = {};
        SETTING_SECTIONS.forEach(function (section) {
            section.fields.forEach(function (field) {
                var node = byId('setting_' + field.key);
                if (!node) {
                    return;
                }
                if (field.type === 'switch') {
                    payload[field.key] = node.checked ? '1' : '0';
                } else if (field.type === 'color') {
                    var text = byId('setting_' + field.key + '_text');
                    payload[field.key] = (text && text.value.trim()) || node.value;
                } else {
                    payload[field.key] = node.value.trim();
                }
            });
        });
        return payload;
    }

    /* --- Site & category forms -------------------------------------------- */

    function categoryOptions(selectedId) {
        return state.categories.map(function (category) {
            var suffix = category.visibility === 1 ? '（登录可见）' : '';
            return '<option value="' + category.id + '"' +
                (category.id === selectedId ? ' selected' : '') + '>' +
                escapeHtml(category.name + suffix) + '</option>';
        }).join('');
    }

    function openSiteForm(site) {
        if (state.categories.length === 0) {
            toast('请先创建一个分类', 'warning');
            switchView('categories');
            return;
        }
        var isEdit = !!site;
        openModal({
            title: isEdit ? '编辑站点' : '新增站点',
            confirmText: isEdit ? '保存修改' : '创建站点',
            bodyHtml:
                '<div class="field"><label for="f_name">站点名称</label>' +
                '<input class="input" id="f_name" value="' + escapeHtml(site ? site.name : '') + '" placeholder="GitHub"></div>' +
                '<div class="field"><label for="f_url">站点地址</label>' +
                '<input class="input" id="f_url" value="' + escapeHtml(site ? site.url : '') + '" placeholder="https://github.com"></div>' +
                '<div class="field"><label for="f_category">所属分类</label>' +
                '<select class="select" id="f_category">' + categoryOptions(site ? site.categoryId : state.categories[0].id) + '</select></div>' +
                '<div class="field"><label for="f_icon">图标地址</label>' +
                '<input class="input" id="f_icon" value="' + escapeHtml(site ? site.icon : '') + '" placeholder="留空则自动使用站点 favicon">' +
                '<span class="hint">图标由服务端下载并缓存到本地，访客浏览器不会直接访问外部站点。</span></div>' +
                '<div class="field"><label for="f_description">站点描述</label>' +
                '<textarea class="textarea" id="f_description">' + escapeHtml(site ? site.description : '') + '</textarea></div>' +
                '<div class="settings-grid">' +
                '<div class="field"><label for="f_sort">排序值</label>' +
                '<input class="input" id="f_sort" type="number" value="' + (site ? site.sort : 0) + '"></div>' +
                '<div class="field"><label>状态</label><label class="switch">' +
                '<input type="checkbox" id="f_status"' + (!site || site.status === 1 ? ' checked' : '') + '>' +
                '<span class="track"></span><span class="label">在导航页显示</span></label></div>' +
                '</div>',
            onConfirm: function () {
                var payload = {
                    name: byId('f_name').value.trim(),
                    url: byId('f_url').value.trim(),
                    categoryId: Number(byId('f_category').value),
                    icon: byId('f_icon').value.trim(),
                    description: byId('f_description').value.trim(),
                    sort: Number(byId('f_sort').value) || 0,
                    status: byId('f_status').checked ? 1 : 0
                };
                if (!payload.name) {
                    toast('站点名称不能为空', 'warning');
                    return false;
                }
                if (!payload.url) {
                    toast('站点地址不能为空', 'warning');
                    return false;
                }
                var options = {
                    method: isEdit ? 'PUT' : 'POST',
                    body: JSON.stringify(isEdit ? Object.assign({ id: site.id }, payload) : payload)
                };
                return request('/admin/site', options).then(function () {
                    toast(isEdit ? '站点已更新' : '站点已创建', 'success');
                    return loadCategories().then(function () {
                        renderSites();
                        renderCategories();
                    });
                }).catch(function (error) {
                    handleError(error);
                    return false;
                });
            }
        });
    }

    function openCategoryForm(category) {
        var isEdit = !!category;
        openModal({
            title: isEdit ? '编辑分类' : '新增分类',
            confirmText: isEdit ? '保存修改' : '创建分类',
            bodyHtml:
                '<div class="field"><label for="c_name">分类名称</label>' +
                '<input class="input" id="c_name" value="' + escapeHtml(category ? category.name : '') + '" placeholder="开发工具"></div>' +
                '<div class="field"><label for="c_icon">分类图标</label>' +
                '<input class="input" id="c_icon" value="' + escapeHtml(category ? category.icon : '') + '" placeholder="可填入 Emoji，例如 🛠️"></div>' +
                '<div class="field"><label for="c_description">分类描述</label>' +
                '<input class="input" id="c_description" value="' + escapeHtml(category ? category.description : '') + '"></div>' +
                '<div class="settings-grid">' +
                '<div class="field"><label for="c_sort">排序值</label>' +
                '<input class="input" id="c_sort" type="number" value="' + (category ? category.sort : 0) + '"></div>' +
                '<div class="field"><label>状态</label><label class="switch">' +
                '<input type="checkbox" id="c_status"' + (!category || category.status === 1 ? ' checked' : '') + '>' +
                '<span class="track"></span><span class="label">在导航页显示</span></label></div>' +
                '</div>' +
                '<div class="field"><label for="c_visibility">可见范围</label>' +
                '<select class="select" id="c_visibility">' +
                '<option value="0"' + (!category || category.visibility !== 1 ? ' selected' : '') + '>公开 · 所有人可见</option>' +
                '<option value="1"' + (category && category.visibility === 1 ? ' selected' : '') + '>登录后可见 · 仅登录管理后台后展示</option>' +
                '</select>' +
                '<span class="hint">设为“登录后可见”后，该分类及其站点不会出现在未登录访客的导航页上。</span></div>',
            onConfirm: function () {
                var payload = {
                    name: byId('c_name').value.trim(),
                    icon: byId('c_icon').value.trim(),
                    description: byId('c_description').value.trim(),
                    sort: Number(byId('c_sort').value) || 0,
                    status: byId('c_status').checked ? 1 : 0,
                    visibility: Number(byId('c_visibility').value) || 0
                };
                if (!payload.name) {
                    toast('分类名称不能为空', 'warning');
                    return false;
                }
                var options = {
                    method: isEdit ? 'PUT' : 'POST',
                    body: JSON.stringify(isEdit ? Object.assign({ id: category.id }, payload) : payload)
                };
                return request('/admin/category', options).then(function () {
                    toast(isEdit ? '分类已更新' : '分类已创建', 'success');
                    return loadCategories().then(function () {
                        renderSites();
                        renderCategories();
                    });
                }).catch(function (error) {
                    handleError(error);
                    return false;
                });
            }
        });
    }

    /* --- Row actions ------------------------------------------------------ */

    function findSite(id) {
        var found = null;
        allSites().forEach(function (row) {
            if (row.site.id === id) {
                found = row.site;
            }
        });
        return found;
    }

    function deleteSite(id) {
        var site = findSite(id);
        if (!site) {
            return;
        }
        confirmDialog('删除站点', '确定要删除「' + site.name + '」吗？该操作不可撤销。', {
            confirmText: '删除',
            danger: true
        }).then(function (confirmed) {
            if (!confirmed) {
                return;
            }
            request('/admin/site?id=' + id, { method: 'DELETE' }).then(function () {
                toast('站点已删除', 'success');
                return loadCategories().then(function () {
                    renderSites();
                    renderCategories();
                });
            }).catch(handleError);
        });
    }

    function deleteCategory(id) {
        var category = state.categories.filter(function (item) {
            return item.id === id;
        })[0];
        if (!category) {
            return;
        }
        confirmDialog('删除分类', '删除「' + category.name + '」会同时删除其下的 ' +
            (category.sites || []).length + ' 个站点，确定继续吗？', {
            confirmText: '全部删除',
            danger: true
        }).then(function (confirmed) {
            if (!confirmed) {
                return;
            }
            request('/admin/category?id=' + id, { method: 'DELETE' }).then(function () {
                toast('分类已删除', 'success');
                return loadCategories().then(function () {
                    renderSites();
                    renderCategories();
                });
            }).catch(handleError);
        });
    }

    /* --- Data & account --------------------------------------------------- */

    function exportData() {
        request('/admin/backup/export').then(function (data) {
            var blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
            var url = URL.createObjectURL(blob);
            var link = document.createElement('a');
            var stamp = new Date().toISOString().slice(0, 19).replace(/[:T]/g, '-');
            link.href = url;
            link.download = 'xyhelper-nav-backup-' + stamp + '.json';
            document.body.appendChild(link);
            link.click();
            link.remove();
            URL.revokeObjectURL(url);
            toast('数据已导出', 'success');
        }).catch(handleError);
    }

    function importData(file) {
        if (!file) {
            return;
        }
        var reader = new FileReader();
        reader.onload = function () {
            var payload;
            try {
                payload = JSON.parse(reader.result);
            } catch (error) {
                toast('文件不是合法的 JSON', 'error');
                return;
            }
            confirmDialog('导入数据', '导入会覆盖当前全部分类与站点，确定继续吗？', {
                confirmText: '覆盖导入',
                danger: true
            }).then(function (confirmed) {
                if (!confirmed) {
                    return;
                }
                request('/admin/backup/import', {
                    method: 'POST',
                    body: JSON.stringify({ data: payload })
                }).then(function () {
                    toast('数据已导入', 'success');
                    return refreshAll();
                }).catch(handleError);
            });
        };
        reader.readAsText(file);
    }

    function refreshAll() {
        return Promise.all([loadCategories(), loadSettings()]).then(function () {
            renderSites();
            renderCategories();
        });
    }

    /* --- Bootstrap -------------------------------------------------------- */

    function enterAdmin() {
        showAdmin();
        return Promise.all([loadCategories(), loadSettings(), loadProfile()]).then(function () {
            renderSites();
            renderCategories();
            switchView(state.view);
        });
    }

    function bindEvents() {
        dom.loginForm.addEventListener('submit', function (event) {
            event.preventDefault();
            var button = dom.loginSubmit;
            button.disabled = true;
            button.textContent = '登录中…';
            dom.loginError.classList.remove('is-visible');
            request('/admin/login', {
                method: 'POST',
                body: JSON.stringify({
                    username: dom.loginUsername.value.trim(),
                    password: dom.loginPassword.value
                })
            }).then(function () {
                return enterAdmin();
            }).catch(function (error) {
                dom.loginError.textContent = error.message || '登录失败';
                dom.loginError.classList.add('is-visible');
            }).then(function () {
                button.disabled = false;
                button.textContent = '登录';
            });
        });

        dom.adminNav.addEventListener('click', function (event) {
            var button = event.target.closest('button[data-view]');
            if (button) {
                switchView(button.getAttribute('data-view'));
            }
        });

        dom.themeToggle.addEventListener('click', function () {
            var next = document.documentElement.getAttribute('data-theme') === 'dark' ? 'light' : 'dark';
            window.localStorage.setItem(THEME_KEY, next);
            applyTheme(next);
        });

        dom.logoutBtn.addEventListener('click', function () {
            request('/admin/logout', { method: 'POST' }).then(function () {
                state.categories = [];
                state.settings = {};
                showLogin();
                dom.loginUsername.value = '';
                dom.loginPassword.value = '';
                toast('已退出登录', 'success');
            }).catch(handleError);
        });

        byId('addSiteBtn').addEventListener('click', function () {
            openSiteForm(null);
        });

        byId('refreshIconsBtn').addEventListener('click', function () {
            var button = this;
            button.disabled = true;
            request('/admin/icons/refresh', {
                method: 'POST',
                body: JSON.stringify({ id: 0 })
            }).then(function (data) {
                var type = data.failed > 0 ? 'warning' : 'success';
                toast('图标刷新完成：成功 ' + data.refreshed + ' 个，失败 ' + data.failed + ' 个', type);
                return loadCategories().then(function () {
                    renderSites();
                    renderCategories();
                });
            }).catch(handleError).then(function () {
                button.disabled = false;
            });
        });

        byId('addCategoryBtn').addEventListener('click', function () {
            openCategoryForm(null);
        });

        dom.siteCategoryFilter.addEventListener('change', function () {
            state.filter.categoryId = Number(dom.siteCategoryFilter.value);
            renderSites();
        });

        dom.siteKeyword.addEventListener('input', function () {
            state.filter.keyword = dom.siteKeyword.value;
            renderSites();
        });

        dom.siteTableBody.addEventListener('click', function (event) {
            var button = event.target.closest('button[data-action]');
            if (!button) {
                return;
            }
            var id = Number(button.getAttribute('data-id'));
            if (button.getAttribute('data-action') === 'edit-site') {
                openSiteForm(findSite(id));
            } else {
                deleteSite(id);
            }
        });

        dom.categoryTableBody.addEventListener('click', function (event) {
            var button = event.target.closest('button[data-action]');
            if (!button) {
                return;
            }
            var id = Number(button.getAttribute('data-id'));
            var category = state.categories.filter(function (item) {
                return item.id === id;
            })[0];
            if (button.getAttribute('data-action') === 'edit-category') {
                openCategoryForm(category);
            } else {
                deleteCategory(id);
            }
        });

        dom.settingsForm.addEventListener('submit', function (event) {
            event.preventDefault();
            request('/admin/settings', {
                method: 'PUT',
                body: JSON.stringify({ settings: collectSettings() })
            }).then(function () {
                toast('设置已保存', 'success');
                return loadSettings();
            }).catch(handleError);
        });

        byId('exportBtn').addEventListener('click', exportData);

        dom.importZone.addEventListener('click', function () {
            dom.importFile.click();
        });

        dom.importFile.addEventListener('change', function () {
            importData(dom.importFile.files[0]);
            dom.importFile.value = '';
        });

        ['dragenter', 'dragover'].forEach(function (type) {
            dom.importZone.addEventListener(type, function (event) {
                event.preventDefault();
                dom.importZone.style.borderColor = 'var(--accent)';
            });
        });

        ['dragleave', 'drop'].forEach(function (type) {
            dom.importZone.addEventListener(type, function (event) {
                event.preventDefault();
                dom.importZone.style.borderColor = '';
            });
        });

        dom.importZone.addEventListener('drop', function (event) {
            if (event.dataTransfer && event.dataTransfer.files.length) {
                importData(event.dataTransfer.files[0]);
            }
        });

        dom.profileForm.addEventListener('submit', function (event) {
            event.preventDefault();
            request('/admin/profile', {
                method: 'PUT',
                body: JSON.stringify({ nickname: dom.profileNickname.value.trim() })
            }).then(function () {
                toast('资料已更新', 'success');
                return loadProfile();
            }).catch(handleError);
        });

        dom.passwordForm.addEventListener('submit', function (event) {
            event.preventDefault();
            request('/admin/password', {
                method: 'POST',
                body: JSON.stringify({
                    oldPassword: dom.oldPassword.value,
                    newPassword: dom.newPassword.value
                })
            }).then(function () {
                toast('密码已更新，请重新登录', 'success');
                dom.oldPassword.value = '';
                dom.newPassword.value = '';
            }).catch(handleError);
        });

        dom.modalConfirm.addEventListener('click', function () {
            var onConfirm = modal.confirm;
            if (!onConfirm) {
                closeModal();
                return;
            }
            var result = onConfirm();
            if (result && typeof result.then === 'function') {
                result.then(function (ok) {
                    if (ok !== false) {
                        closeModal();
                    }
                });
            } else if (result !== false) {
                closeModal();
            }
        });

        dom.modalCancel.addEventListener('click', cancelModal);
        dom.modalClose.addEventListener('click', cancelModal);
        dom.modalMask.addEventListener('click', function (event) {
            if (event.target === dom.modalMask) {
                cancelModal();
            }
        });
        document.addEventListener('keydown', function (event) {
            if (event.key === 'Escape' && dom.modalMask.classList.contains('is-open')) {
                cancelModal();
            }
        });
    }

    function init() {
        dom = {
            loginView: byId('loginView'),
            loginForm: byId('loginForm'),
            loginError: byId('loginError'),
            loginUsername: byId('loginUsername'),
            loginPassword: byId('loginPassword'),
            loginSubmit: byId('loginSubmit'),
            adminView: byId('adminView'),
            adminNav: byId('adminNav'),
            viewTitle: byId('viewTitle'),
            viewSubtitle: byId('viewSubtitle'),
            themeToggle: byId('themeToggle'),
            logoutBtn: byId('logoutBtn'),
            sideName: byId('sideName'),
            sideLogo: byId('sideLogo'),
            statGrid: byId('statGrid'),
            siteCategoryFilter: byId('siteCategoryFilter'),
            siteKeyword: byId('siteKeyword'),
            siteTableBody: byId('siteTableBody'),
            categoryTableBody: byId('categoryTableBody'),
            settingsForm: byId('settingsForm'),
            importZone: byId('importZone'),
            importFile: byId('importFile'),
            profileForm: byId('profileForm'),
            profileNickname: byId('profileNickname'),
            passwordForm: byId('passwordForm'),
            oldPassword: byId('oldPassword'),
            newPassword: byId('newPassword'),
            accountMeta: byId('accountMeta'),
            modalMask: byId('modalMask'),
            modalTitle: byId('modalTitle'),
            modalBody: byId('modalBody'),
            modalConfirm: byId('modalConfirm'),
            modalCancel: byId('modalCancel'),
            modalClose: byId('modalClose'),
            toastHost: byId('toastHost')
        };

        applyTheme(initialTheme());
        bindEvents();

        request('/admin/profile').then(function (data) {
            state.profile = data;
            return enterAdmin();
        }).catch(function () {
            showLogin();
            dom.loginUsername.focus();
        });
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
