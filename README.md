# easy-menu-go

```typescript
type CommandSpec =
  | string
  | {
      run: string;
      env: { [key: string]: string };
      work_dir: string;
    };

type Command = {
  [commandName: string]: CommandSpec;
};
type EvalMenu = {
  eval: string;
};

type Menu = {
  [menuName: string]: Array<Command | Menu | EvalMenu>;
  env: { [key: string]: string };
  work_dir: string;
};
```
