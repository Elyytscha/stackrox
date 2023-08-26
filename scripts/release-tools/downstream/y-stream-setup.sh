#!/usr/bin/env bash
set -eoux pipefail

# Check the number of input parameters
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <previous_release> <release>"
    exit 1
fi

previous_release="$1"
release="$2"
ship_date="$3"

# Define the version pattern
version_pattern="^[0-9]+\.[0-9]+$"

# Check if variables match the version pattern
if ! [[ $previous_release =~ $version_pattern ]] || ! [[ $release =~ $version_pattern ]]; then
    echo "Both previous_release and release must match the pattern x.y, where x and y are integers."
    exit 1
fi

# Extract major and minor version numbers
previous_major="$(cut -d '.' -f 1 <<< "$previous_release")"
previous_minor="$(cut -d '.' -f 2 <<< "$previous_release")"
release_major="$(cut -d '.' -f 1 <<< "$release")"
release_minor="$(cut -d '.' -f 2 <<< "$release")"

# Compare versions
if [ "$previous_major" -gt "$release_major" ] || { [ "$previous_major" -eq "$release_major" ] && [ "$previous_minor" -gt "$release_minor" ] ; }; then
    echo "Previous release must be less than the current release."
    exit 1
fi

# Define the ship_date pattern
ship_date_pattern="^[0-9]{4}\-[0-9]{2}\-[0-9]{2}$"

if ! [[ $ship_date =~ $ship_date_pattern ]]; then
    echo "Ship date is not valid"
fi

get_supported_versions() {

    mapfile -t supported_versions < <(
      curl -fsSL "https://access.redhat.com/product-life-cycles/api/v1/products?name=Red%20Hat%20Advanced%20Cluster%20Security%20for%20Kubernetes" |
      jq -r '.data[0].versions[] | select(.type == "Full Support") | .name'
    )


    nversions=${#supported_versions[@]}
    for ((i = nversions - 1; i >= 0; i = i - 1)); do
        echo "${supported_versions[$i]}"
    done

    echo "$release"
}

update_content_stream_tags() {
    mapfile -t versions < <(get_supported_versions)

    nversions=${#versions[@]}

    for ((i=0; i<nversions; i++)); do
        find versions/release-* -name 'product.yml' -exec bash -c 'yq w -i "$1" delivery-repo-content.content_stream_tags[$2] "$3" --style=double' _ {} "$i" "${versions[$i]}" \;
    done
}

perform_sed_or_gsed() {
    local search="$1"
    local replace="$2"
    local file="$3"

    if command -v sed > /dev/null; then
        sed -i "s|$search|$replace|" "$file"
    elif command -v gsed > /dev/null; then
        gsed -i "s|$search|$replace|" "$file"
    else
        echo "Error: Neither sed nor gsed found. Cannot perform replacement."
        exit 1
    fi
}

export -f perform_sed_or_gsed

git clone git@gitlab.cee.redhat.com:cpaas-products/rhacs.git gitlab-rhacs
pushd gitlab-rhacs

git checkout -b "setup-${release}"

cp -R "versions/release-${previous_release}" "versions/release-${release}"

pushd "versions/release-${release}"

rm -f advisory_map.yml

# Update release.yml with the correct version. Sinple sed should be safe, but check changes
perform_sed_or_gsed "$previous_release" "$release" release.yml

# Update product.yml with the correct version.
# This is a little more complicated since not all occurances of the old release should be changed
# Check changes
perform_sed_or_gsed "rhacs-$previous_release" "rhacs-$release" product.yml
perform_sed_or_gsed "RHACS-$previous_release" "RHACS-$release" product.yml
perform_sed_or_gsed "RHACS $previous_release" "RHACS $release" product.yml
perform_sed_or_gsed "Kubernetes $previous_release" "Kubernetes $release" product.yml

yq w -i product.yml product.release.version "${release}.0" --style=single
yq w -i product.yml honeybadger.version "${release}" --style=single
yq w -i product.yml product.honeybadger.version "${release}" --style=single

# Update the ship_date
perform_sed_or_gsed 'ship_date: *' 'ship_date: "$ship_date"' product.yml

popd

# Add the versions to content_stream_tags in all product.yml files
update_content_stream_tags

# yq makes some unwanted changes that need to be undone
find versions/release-* -name 'product.yml' -exec bash -c 'perform_sed_or_gsed "!!merge " "" "$1"' _ {} \;

# Add --- to the beginning of all product.yml files which is removed by yq
find versions/release-* -name 'product.yml' -exec bash -c 'echo "---" > temp.txt; cat "$1" >> temp.txt; mv temp.txt "$1"' _ {} \;

echo
echo
echo
echo

# TODO Automate this once there is confidance that this script is working.
echo "ATTENTION: Manually check, commit, and push the changes. Create an MR."
