---
id: def-new-types
title: Defining New Types
sidebar_label: Defining New Types
sidebar_position: 39
hide_title: true
draft: false
---

## Introduction

We can define new data types with the help of `type` directive:

```
type <Name> <definition>
```

In its simplest form is possible to define a new type from an existing one. For example,

```
type Age uint8
```

The new `Age` type is actually a `uint8` that accepts values ​​between `0` and `255` and allows you to perform the same operations as a `uint8`. So why not directly use a uint8 to save the data pertaining to age? The main reason is that the Age data type brings a new semantic meaning within the program domain: it is a number, but we know that it only contains age-related data.



