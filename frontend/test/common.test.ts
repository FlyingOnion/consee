import { expect, test } from "bun:test";
import { addNewKey, parseKey, TreeData, TreeItem } from "../src/common/kz";

test("parse key", () => {
  let dataList = parseKey("a", "/");
  expect<TreeData[]>(dataList).toContainEqual({
    key: "a",
    path: "a",
    depth: 1,
    isLeaf: true,
  });

  dataList = parseKey("a/", "/");
  expect<TreeData[]>(dataList).toContainEqual({
    key: "a",
    path: "a/",
    depth: 1,
    isLeaf: false,
  });

  dataList = parseKey("a/b/c", "/");
  expect<TreeData[]>(dataList).toContainEqual({
    key: "a",
    path: "a/",
    depth: 1,
    isLeaf: false,
  });
  expect<TreeData[]>(dataList).toContainEqual({
    key: "b",
    path: "a/b/",
    depth: 2,
    isLeaf: false,
  });
  expect<TreeData[]>(dataList).toContainEqual({
    key: "c",
    path: "a/b/c",
    depth: 3,
    isLeaf: true,
  });
});

test("add new key", () => {
  const tree: TreeItem[] = [];
  addNewKey(tree, "a");
  expect(tree).toEqual([
    {
      key: "a",
      path: "a",
      depth: 1,
      isLeaf: true,
      children: [],
      access: "",
    },
  ]);

  addNewKey(tree, "a/");
  expect(tree).toEqual([
    {
      key: "a",
      path: "a",
      depth: 1,
      isLeaf: true,
      children: [],
      access: "",
    },
    {
      key: "a",
      path: "a/",
      depth: 1,
      isLeaf: false,
      children: [],
      access: "",
    },
  ]);

  addNewKey(tree, "a/b/c");
  expect(tree).toEqual([
    {
      key: "a",
      path: "a",
      depth: 1,
      isLeaf: true,
      children: [],
      access: "",
    },
    {
      key: "a",
      path: "a/",
      depth: 1,
      isLeaf: false,
      children: [
        {
          key: "b",
          path: "a/b/",
          depth: 2,
          isLeaf: false,
          children: [
            {
              key: "c",
              path: "a/b/c",
              depth: 3,
              isLeaf: true,
              children: [],
              access: "",
            },
          ],
          access: "",
        },
      ],
      access: "",
    },
  ]);
});
