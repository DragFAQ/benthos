---
title: azure_blob_storage
type: output
status: beta
categories: ["Services","Azure"]
---

<!--
     THIS FILE IS AUTOGENERATED!

     To make changes please edit the contents of:
     lib/output/azure_blob_storage.go
-->

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

:::caution BETA
This component is mostly stable but breaking changes could still be made outside of major version releases if a fundamental problem with the component is found.
:::

Sends message parts as objects to an Azure Blob Storage Account container. Each
object is uploaded with the filename specified with the `container`
field.

Introduced in version 3.36.0.


<Tabs defaultValue="common" values={[
  { label: 'Common', value: 'common', },
  { label: 'Advanced', value: 'advanced', },
]}>

<TabItem value="common">

```yml
# Common config fields, showing default values
output:
  label: ""
  azure_blob_storage:
    storage_account: ""
    storage_access_key: ""
    storage_sas_token: ""
    storage_connection_string: ""
    container: ""
    path: ${!count("files")}-${!timestamp_unix_nano()}.txt
    max_in_flight: 64
```

</TabItem>
<TabItem value="advanced">

```yml
# All config fields, showing default values
output:
  label: ""
  azure_blob_storage:
    storage_account: ""
    storage_access_key: ""
    storage_sas_token: ""
    storage_connection_string: ""
    public_access_level: PRIVATE
    container: ""
    path: ${!count("files")}-${!timestamp_unix_nano()}.txt
    blob_type: BLOCK
    max_in_flight: 64
```

</TabItem>
</Tabs>

Only one authentication method is required, `storage_connection_string` or `storage_account` and `storage_access_key`. If both are set then the `storage_connection_string` is given priority.

In order to have a different path for each object you should use function
interpolations described [here](/docs/configuration/interpolation#bloblang-queries), which are
calculated per message of a batch.

## Performance

This output benefits from sending multiple messages in flight in parallel for
improved performance. You can tune the max number of in flight messages with the
field `max_in_flight`.

## Fields

### `storage_account`

The storage account to upload messages to. This field is ignored if `storage_connection_string` is set.


Type: `string`  
Default: `""`  

### `storage_access_key`

The storage account access key. This field is ignored if `storage_connection_string` is set.


Type: `string`  
Default: `""`  

### `storage_sas_token`

The storage account SAS token. This field is ignored if `storage_connection_string` or `storage_access_key` / `storage_sas_token` are set.


Type: `string`  
Default: `""`  
Requires version 3.38.0 or newer  

### `storage_connection_string`

A storage account connection string. This field is required if `storage_account` and `storage_access_key` are not set.


Type: `string`  
Default: `""`  

### `public_access_level`

The container's public access level. The default value is `PRIVATE`.


Type: `string`  
Default: `"PRIVATE"`  
Options: `PRIVATE`, `BLOB`, `CONTAINER`.

### `container`

The container for uploading the messages to.
This field supports [interpolation functions](/docs/configuration/interpolation#bloblang-queries).


Type: `string`  
Default: `""`  

```yml
# Examples

container: messages-${!timestamp("2006")}
```

### `path`

The path of each message to upload.
This field supports [interpolation functions](/docs/configuration/interpolation#bloblang-queries).


Type: `string`  
Default: `"${!count(\"files\")}-${!timestamp_unix_nano()}.txt"`  

```yml
# Examples

path: ${!count("files")}-${!timestamp_unix_nano()}.json

path: ${!meta("kafka_key")}.json

path: ${!json("doc.namespace")}/${!json("doc.id")}.json
```

### `blob_type`

Block and Append blobs are comprised of blocks, and each blob can support up to 50,000 blocks. The default value is `+"`BLOCK`"+`.`
This field supports [interpolation functions](/docs/configuration/interpolation#bloblang-queries).


Type: `string`  
Default: `"BLOCK"`  
Options: `BLOCK`, `APPEND`.

### `max_in_flight`

The maximum number of messages to have in flight at a given time. Increase this to improve throughput.


Type: `int`  
Default: `64`  


