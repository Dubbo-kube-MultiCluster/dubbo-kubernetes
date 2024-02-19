/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { ref } from 'vue'

export const PRIMARY_COLOR_DEFAULT = '#17b392'

export const LOCAL_STORAGE_LOCALE = 'LOCAL_STORAGE_LOCALE'
export const LOCAL_STORAGE_THEME = 'LOCAL_STORAGE_THEME'

let item = localStorage.getItem(LOCAL_STORAGE_THEME)

export const PRIMARY_COLOR = ref(item || PRIMARY_COLOR_DEFAULT)
