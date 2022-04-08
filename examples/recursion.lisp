;; this is recursive
(define factorial
  (lambda (n)
    (if (= n 0)
        1
        (* n (factorial (- n 1))))))

;; this is tail-recursive/iterative
;; except this Lisp interpreter doesn't implement tail recursion ¯\_(ツ)_/¯
(define factorial-iter
  (lambda (n)
    (define do-factorial-iter
      (lambda (n acc)
        (if (= n 0)
            acc
            (do-factorial-iter (- n 1) (* acc n)))))
    (do-factorial-iter n 1)))

(print 'recursive)
(factorial 1)
(factorial 2)
(factorial 3)
(factorial 4)
(factorial 5)

(print 'iterative)
(factorial-iter 1)
(factorial-iter 2)
(factorial-iter 3)
(factorial-iter 4)
(factorial-iter 5)
