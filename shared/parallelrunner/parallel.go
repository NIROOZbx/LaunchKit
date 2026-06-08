package parallelrunner

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func Query2[A, B any](ctx context.Context,
	fnA func( context.Context)(A,error),
	fnB func( context.Context)(B,error),
	) (A,B,error) {
		var a A
		var b B
		
		g,ctx:=errgroup.WithContext(ctx)

		g.Go(func() error{
			var err error
			a, err = fnA(ctx)
			return err
		})

		
		g.Go(func() error{
			var err error
			b, err = fnB(ctx)
			return err
		})

		if err:=g.Wait();err !=nil{
			return a,b,err
		}
		return a,b,nil
}


func Query3[A, B, C any](ctx context.Context,
	fnA func(context.Context) (A, error),
	fnB func(context.Context) (B, error),
	fnC func(context.Context) (C, error),
) (A, B, C, error) {
	var a A
	var b B
	var c C

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		a, err = fnA(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		b, err = fnB(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		c, err = fnC(ctx)
		return err
	})

	if err := g.Wait(); err != nil {
		return a, b, c, err
	}

	return a, b, c, nil
}

func Query4[A, B, C, D any](ctx context.Context,
	fnA func(context.Context) (A, error),
	fnB func(context.Context) (B, error),
	fnC func(context.Context) (C, error),
	fnD func(context.Context) (D, error),
) (A, B, C, D, error) {
	var a A
	var b B
	var c C
	var d D

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		a, err = fnA(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		b, err = fnB(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		c, err = fnC(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		d, err = fnD(ctx)
		return err
	})

	if err := g.Wait(); err != nil {
		return a, b, c, d, err
	}

	return a, b, c, d, nil
}

func Query5[A, B, C, D, E any](ctx context.Context,
	fnA func(context.Context) (A, error),
	fnB func(context.Context) (B, error),
	fnC func(context.Context) (C, error),
	fnD func(context.Context) (D, error),
	fnE func(context.Context) (E, error),
) (A, B, C, D, E, error) {
	var a A
	var b B
	var c C
	var d D
	var e E

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		a, err = fnA(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		b, err = fnB(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		c, err = fnC(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		d, err = fnD(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		e, err = fnE(ctx)
		return err
	})

	if err := g.Wait(); err != nil {
		return a, b, c, d, e, err
	}

	return a, b, c, d, e, nil
}