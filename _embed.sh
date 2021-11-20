#!/bin/bash
# Copyright © 2021 Kris Nóva <kris@nivenly.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
#   ███╗   ██╗ █████╗ ███╗   ███╗██╗
#   ████╗  ██║██╔══██╗████╗ ████║██║
#   ██╔██╗ ██║███████║██╔████╔██║██║
#   ██║╚██╗██║██╔══██║██║╚██╔╝██║██║
#   ██║ ╚████║██║  ██║██║ ╚═╝ ██║███████╗
#   ╚═╝  ╚═══╝╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝
#

main=$(cat ./src/main.go.tpl)
library=$(cat ./src/library.go.tpl)

rm -f embed_gen.go
touch embed_gen.go

cat > embed_gen.go << EOL
//
// Copyright © 2021 Kris Nóva <kris@nivenly.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//   ███╗   ██╗ █████╗ ███╗   ███╗██╗
//   ████╗  ██║██╔══██╗████╗ ████║██║
//   ██╔██╗ ██║███████║██╔████╔██║██║
//   ██║╚██╗██║██╔══██║██║╚██╔╝██║██║
//   ██║ ╚████║██║  ██║██║ ╚═╝ ██║███████╗
//   ╚═╝  ╚═══╝╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝
//

EOL

echo "package naml" >> embed_gen.go
echo "" >> embed_gen.go
echo "const FormatMainGo string = \`" >> embed_gen.go
echo "$main" >> embed_gen.go
echo "\`" >> embed_gen.go
echo "" >> embed_gen.go
echo "const FormatLibraryGo string = \`" >> embed_gen.go
echo "$library" >> embed_gen.go
echo "\`" >> embed_gen.go
echo "" >> embed_gen.go

