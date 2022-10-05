---
id: index_templates
title: Index Templates
sidebar_label: Index Templates
sidebar_position: 19
---

An index template specify how to configure an index when it is created. For data streams, the index template configures the stream's backing indices as they are created. Templates are configured prior to index creation. When index is created - either manually or through indexing a document - the template settings are used as a basis for creating the index. 

There are two types of templates: index templates and component templates. Component templates are reusable building blocks that configure mappings, settings, and aliases. While you can use component templates to construct index templates, they aren’t directly applied to a set of indices. Index templates can contain a collection of component templates, as well as directly specify settings, mappings, and aliases.

The following conditions apply to index templates:

+ Composable templates take precedence over legacy templates. If no composable template matches a given index, a legacy template may still match and be applied.
+ If an index is created with explicit settings and also matches an index template, the settings from the create index request take precedence over settings specified in the index template and its component templates.
+ If a new data stream or index matches more than one index template, the index template with the highest priority is used.